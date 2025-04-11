package main

import (
	"context"
	_ "expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/facade"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/authapi"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/redis"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	conf "github.com/OutOfStack/game-library/internal/pkg/config"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	zaplog "github.com/OutOfStack/game-library/internal/pkg/log"
	"github.com/OutOfStack/game-library/internal/taskprocessor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

// @title Game library API
// @version 0.3
// @description API for game library service
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /api
// @query.collection.format multi
// @schemes http
func main() {
	// load config
	var cfg appconf.Cfg
	if err := conf.Load(".", "app", "env", &cfg); err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	// init logger
	logger := zaplog.New(cfg)
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			log.Printf("can't sync logger: %v", err)
		}
	}(logger)

	// run
	if err := run(logger, cfg); err != nil {
		logger.Fatal("can't run app", zap.Error(err))
	}
}

func run(logger *zap.Logger, cfg appconf.Cfg) error {
	ctx := context.Background()

	// connect to database
	db, err := database.New(ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("connect to db: %v", err)
	}
	defer db.Close()

	// create IGDB client
	igdbAPIClient, err := igdbapi.New(logger, cfg.IGDB)
	if err != nil {
		return fmt.Errorf("create IGDB client: %w", err)
	}

	// create auth api client
	authAPIClient, err := authapi.New(logger, cfg.Auth.VerifyTokenAPIURL)
	if err != nil {
		return fmt.Errorf("create auth api client: %w", err)
	}

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

	// create redis cache service
	rCache := cache.NewRedisStore(redisClient, logger)

	// create storage
	storage := repo.New(db)

	// create auth facade
	authFacade, err := auth.New(logger, cfg.Auth.SigningAlgorithm, authAPIClient)
	if err != nil {
		return fmt.Errorf("create Auth: %w", err)
	}

	// create game facade
	gameFacade := facade.NewProvider(logger, storage, rCache)

	// create api provider
	apiProvider := api.NewProvider(logger, rCache, gameFacade)

	// run background tasks
	taskProvider := taskprocessor.New(logger, storage, igdbAPIClient, s3Client)
	scheduler := gocron.NewScheduler(time.UTC)
	tasks := map[string]model.TaskInfo{
		taskprocessor.FetchIGDBGamesTaskName: {Schedule: cfg.Scheduler.FetchIGDBGames, Fn: taskProvider.StartFetchIGDBGames},
	}
	for name, task := range tasks {
		_, err = scheduler.Cron(task.Schedule).Name(name).Do(task.Fn)
		if err != nil {
			logger.Error("run task", zap.String("task", name), zap.Error(err))
			return fmt.Errorf("run task %s: %v", name, err)
		}
	}
	scheduler.StartAsync()

	// start debug service
	go func() {
		logger.Info("Debug service started", zap.String("address", cfg.Web.DebugAddress))
		profilerRouter := chi.NewRouter()
		profilerRouter.Mount("/debug", middleware.Profiler())
		debugService := http.Server{
			Addr:        cfg.Web.DebugAddress,
			Handler:     profilerRouter,
			ReadTimeout: time.Second,
		}
		err = debugService.ListenAndServe()
		if err != nil {
			logger.Error("Debug service stopped", zap.Error(err))
		}
	}()

	// start API service
	apiService, err := api.Service(logger, db, authFacade, apiProvider, cfg)
	if err != nil {
		return fmt.Errorf("can't create service api: %w", err)
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("API service started", zap.String("address", apiService.Addr))
		serverErrors <- apiService.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err = <-serverErrors:
		return fmt.Errorf("listening and serving: %w", err)
	case <-shutdown:
		logger.Info("Start shutdown")
		timeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err = apiService.Shutdown(ctx); err != nil {
			logger.Error("Shutdown did not complete", zap.Duration("timeout", timeout), zap.Error(err))
			err = apiService.Close()
			if err != nil {
				return fmt.Errorf("shutdown: %w", err)
			}
		}
	}

	return nil
}
