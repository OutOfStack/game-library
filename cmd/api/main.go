package main

import (
	"context"
	_ "expvar"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/api"
	"github.com/OutOfStack/game-library/internal/api/grpc/infoapi"
	"github.com/OutOfStack/game-library/internal/api/tools"
	"github.com/OutOfStack/game-library/internal/api/validation"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/authapi"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/redis"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/facade"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/pkg/logging"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/OutOfStack/game-library/internal/repo"
	infopb "github.com/OutOfStack/game-library/pkg/proto/infoapi/v1"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpcrefl "google.golang.org/grpc/reflection"
)

func main() {
	// load config
	cfg, err := appconf.Load[appconf.API]()
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

func run(logger *zap.Logger, cfg *appconf.API) error {
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
	gameFacade := facade.NewProvider(logger, storage, cacheStore, s3Client, nil, igdbAPIClient)

	// create request validation rules
	validationRules := validation.NewRules(logger, cfg)

	// create api provider
	apiProvider := api.NewProvider(logger, cacheStore, gameFacade, validationRules)

	serverErrors := make(chan error, 3)

	// start debug service; it serves probes, so it starts first
	hc := tools.NewHealthCheck(db)
	debugService := observability.NewDebugServer(cfg.Debug.Address, hc.Liveness, hc.Readiness)
	go func() {
		logger.Info("Debug service started", zap.String("address", cfg.Debug.Address))
		serverErrors <- debugService.ListenAndServe()
	}()

	// start http API service
	apiService, tracerProvider, err := api.Service(logger, authFacade, apiProvider, cfg)
	if err != nil {
		return fmt.Errorf("can't create service api: %w", err)
	}
	apiListenConfig := net.ListenConfig{}
	apiListener, err := apiListenConfig.Listen(ctx, "tcp", apiService.Addr)
	if err != nil {
		return fmt.Errorf("failed to create http listener: %w", err)
	}
	go func() {
		logger.Info("API service started", zap.String("address", apiService.Addr))
		serverErrors <- apiService.Serve(apiListener)
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

	hc.SetReady(true)

	// wait for shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))

		hc.SetReady(false)
		time.Sleep(5 * time.Second)

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
