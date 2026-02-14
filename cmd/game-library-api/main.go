package main

import (
	"context"
	_ "expvar"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/api"
	"github.com/OutOfStack/game-library/internal/api/grpc/infoapi"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/authapi"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/openaiapi"
	"github.com/OutOfStack/game-library/internal/client/redis"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/facade"
	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/pkg/logging"
	"github.com/OutOfStack/game-library/internal/repo"
	"github.com/OutOfStack/game-library/internal/taskprocessor"
	"github.com/OutOfStack/game-library/internal/web"
	infopb "github.com/OutOfStack/game-library/pkg/proto/infoapi/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpcrefl "google.golang.org/grpc/reflection"
)

// @title Game library API
// @version 0.4
// @description API for game library service
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /api
// @query.collection.format multi
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth
func main() {
	// load config
	cfg, err := appconf.Get()
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	// init logger
	logger := logging.New(cfg.Log.Level, cfg.Graylog.Address)
	defer func() {
		if sErr := logger.Sync(); sErr != nil {
			logger.Error("can't sync logger: %v", zap.Error(sErr))
		}
	}()

	// run
	if err = run(logger, cfg); err != nil {
		logger.Fatal("can't run app", zap.Error(err))
	}
}

func run(logger *zap.Logger, cfg *appconf.Cfg) error {
	ctx := context.Background()

	// connect to database
	db, err := database.New(ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("connect to db: %v", err)
	}
	defer db.Close()

	// create redis client
	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return fmt.Errorf("create redis client: %w", err)
	}

	// create s3 client
	s3Client, err := s3.New(logger, cfg.S3)
	if err != nil {
		return fmt.Errorf("create S3 client: %w", err)
	}

	// create IGDB client
	igdbAPIClient, err := igdbapi.New(logger, cfg.IGDB)
	if err != nil {
		return fmt.Errorf("create IGDB client: %w", err)
	}

	// create authapi client
	grpcClientMetrics := grpcprom.NewClientMetrics()
	err = prometheus.Register(grpcClientMetrics)
	if err != nil {
		return fmt.Errorf("register prometheus grpc client metrics: %w", err)
	}
	authAPIClient, err := authapi.NewClient(authapi.Config{
		Address: cfg.AuthAPI.Address,
		Timeout: cfg.AuthAPI.Timeout,
		DialOptions: []grpc.DialOption{
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithChainUnaryInterceptor(grpcClientMetrics.UnaryClientInterceptor()),
		},
	})
	if err != nil {
		return fmt.Errorf("create auth api client: %w", err)
	}
	defer func() {
		if err = authAPIClient.Close(); err != nil {
			logger.Error("close auth api client", zap.Error(err))
		}
	}()

	// create openai client
	openAIClient := openaiapi.New(logger, cfg.OpenAI)

	// create redis cache service
	cacheStore := cache.NewRedisStore(redisClient, logger)

	// create storage
	storage := repo.New(db, logger)

	// create auth facade
	authFacade, err := auth.New(logger, authAPIClient)
	if err != nil {
		return fmt.Errorf("create Auth: %w", err)
	}

	// create game facade
	gameFacade := facade.NewProvider(logger, storage, cacheStore, s3Client, openAIClient, igdbAPIClient)

	// create web decoder
	decoder := web.NewDecoder(logger, cfg)

	// create api provider
	apiProvider := api.NewProvider(logger, cacheStore, gameFacade, decoder)

	// run background tasks
	taskProvider := taskprocessor.New(logger, storage, igdbAPIClient, s3Client, gameFacade, gameFacade)
	scheduler := gocron.NewScheduler(time.UTC)
	tasks := map[string]model.TaskInfo{
		taskprocessor.FetchIGDBGamesTaskName:      {Schedule: cfg.Scheduler.FetchIGDBGames, Fn: taskProvider.StartFetchIGDBGames},
		taskprocessor.UpdateTrendingIndexTaskName: {Schedule: cfg.Scheduler.UpdateTrendingIndex, Fn: taskProvider.StartUpdateTrendingIndex},
		taskprocessor.UpdateGameInfoTaskName:      {Schedule: cfg.Scheduler.UpdateGameInfo, Fn: taskProvider.StartUpdateGameInfo},
		taskprocessor.ProcessModerationTaskName:   {Schedule: cfg.Scheduler.ProcessModeration, Fn: taskProvider.StartProcessModeration},
	}
	for name, task := range tasks {
		_, err = scheduler.Cron(task.Schedule).Name(name).Do(task.Fn)
		if err != nil {
			logger.Error("run task", zap.String("task", name), zap.Error(err))
			return fmt.Errorf("run task %s: %v", name, err)
		}
	}
	scheduler.StartAsync()

	serverErrors := make(chan error, 3)

	// start debug service
	profilerRouter := chi.NewRouter()
	profilerRouter.Mount("/debug", middleware.Profiler())
	debugService := http.Server{
		Addr:        cfg.Web.DebugAddress,
		Handler:     profilerRouter,
		ReadTimeout: time.Second,
	}
	go func() {
		logger.Info("Debug service started", zap.String("address", cfg.Web.DebugAddress))
		serverErrors <- debugService.ListenAndServe()
	}()

	// start http API service
	apiService, tracerProvider, err := api.Service(logger, db, authFacade, apiProvider, cfg)
	if err != nil {
		return fmt.Errorf("can't create service api: %w", err)
	}
	go func() {
		logger.Info("API service started", zap.String("address", apiService.Addr))
		serverErrors <- apiService.ListenAndServe()
	}()

	// start grpc service
	grpcServerMetrics := grpcprom.NewServerMetrics()
	err = prometheus.Register(grpcServerMetrics)
	if err != nil {
		return fmt.Errorf("register prometheus grpc server metrics: %w", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(grpcServerMetrics.UnaryServerInterceptor()),
	)
	infoService := infoapi.NewInfoService(logger, gameFacade)
	infopb.RegisterInfoApiServiceServer(grpcServer, infoService)
	// register reflection service for grpcurl and other tools
	grpcrefl.Register(grpcServer)
	grpcServerMetrics.InitializeMetrics(grpcServer)

	grpcListenConfig := net.ListenConfig{}
	listener, err := grpcListenConfig.Listen(ctx, "tcp", cfg.Web.GRPCAddress)
	if err != nil {
		return fmt.Errorf("failed to create gRPC listener: %w", err)
	}
	go func() {
		logger.Info("gRPC service started", zap.String("address", cfg.Web.GRPCAddress))
		serverErrors <- grpcServer.Serve(listener)
	}()

	// wait for shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))

		bCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var wg sync.WaitGroup

		// stop grpc server
		wg.Go(func() {
			grpcServer.GracefulStop()
		})

		// stop http server
		wg.Go(func() {
			if shutdownErr := apiService.Shutdown(bCtx); shutdownErr != nil {
				logger.Error("http service shutdown failed", zap.Error(shutdownErr))
				if shutdownErr = apiService.Close(); shutdownErr != nil {
					logger.Error("http service force shutdown failed", zap.Error(shutdownErr))
				}
			}
		})

		// stop debug service
		wg.Go(func() {
			logger.Info("stop debug service")
			if shutdownErr := debugService.Shutdown(bCtx); shutdownErr != nil {
				logger.Error("debug service shutdown failed", zap.Error(shutdownErr))
				if shutdownErr = debugService.Close(); shutdownErr != nil {
					logger.Error("debug service force shutdown failed", zap.Error(shutdownErr))
				}
			}
		})

		// stop tracer provider
		wg.Go(func() {
			logger.Info("stop tracer provider")
			if shutdownErr := tracerProvider.Shutdown(bCtx); shutdownErr != nil {
				logger.Error("tracer provider shutdown failed", zap.Error(shutdownErr))
			}
		})

		wg.Wait()
	}

	return nil
}
