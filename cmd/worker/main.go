package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/openaiapi"
	"github.com/OutOfStack/game-library/internal/client/redis"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/facade"
	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/pkg/logging"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/OutOfStack/game-library/internal/repo"
	"github.com/OutOfStack/game-library/internal/taskprocessor"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

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
		logger.Fatal("can't run worker", zap.Error(err))
	}
}

func run(logger *zap.Logger, cfg *appconf.Cfg) error {
	ctx := context.Background()
	tracerProvider, err := observability.InitTracer(logger, cfg.Jaeger.OTLPEndpoint, "game-library-worker")
	if err != nil {
		return fmt.Errorf("initialize tracer: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if shutdownErr := tracerProvider.Shutdown(shutdownCtx); shutdownErr != nil {
			logger.Error("tracer shutdown failed", zap.Error(shutdownErr))
		}
	}()

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

	// create openai client
	openAIClient := openaiapi.New(logger, cfg.OpenAI)

	// create redis cache service
	cacheStore := cache.NewRedisStore(redisClient, logger)

	// create storage
	storage := repo.New(db, logger)

	// create game facade
	gameFacade := facade.NewProvider(logger, storage, cacheStore, s3Client, openAIClient, igdbAPIClient)

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
			return fmt.Errorf("run task %s: %v", name, err)
		}
	}
	// bind before starting jobs so a metrics port conflict fails startup
	metricsServer := observability.NewMetricsServer(cfg.Web.HTTPAddress)
	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", metricsServer.Addr)
	if err != nil {
		return fmt.Errorf("listen for metrics: %w", err)
	}
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- metricsServer.Serve(listener)
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if shutdownErr := metricsServer.Shutdown(shutdownCtx); shutdownErr != nil {
			logger.Error("metrics shutdown failed", zap.Error(shutdownErr))
			if closeErr := metricsServer.Close(); closeErr != nil {
				logger.Error("metrics close failed", zap.Error(closeErr))
			}
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdown)
	scheduler.StartAsync()
	defer scheduler.Stop()
	logger.Info("Worker started")

	select {
	case sig := <-shutdown:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	case err = <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("metrics server: %w", err)
		}
	}

	return nil
}
