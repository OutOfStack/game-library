package main

import (
	"context"
	_ "expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/handler"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/appconf"
	auth_ "github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"github.com/OutOfStack/game-library/internal/client/uploadcare"
	conf "github.com/OutOfStack/game-library/internal/pkg/config"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/taskprocessor"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title Game library API
// @version 0.2
// @description API for game library service
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /api
// @query.collection.format multi
// @schemes http
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.DisableCaller = true
	loggerCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Printf("syncing logger: %v", err)
		}
	}(logger)

	var cfg appconf.Cfg
	if err = conf.Load(".", "app", "env", &cfg); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		RequireSSL: cfg.DB.RequireSSL,
	})
	if err != nil {
		return fmt.Errorf("opening db: %w", err)
	}
	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			logger.Error("calling database close", zap.Error(err))
		}
	}(db)

	// create auth module
	auth, err := auth_.New(logger, cfg.Auth.SigningAlgorithm, cfg.Auth.VerifyTokenAPIURL)
	if err != nil {
		return fmt.Errorf("creating Auth: %w", err)
	}

	// create IGDB client
	igdbClient, err := igdb.New(logger, cfg.IGDB)
	if err != nil {
		return fmt.Errorf("creating IGDB client: %w", err)
	}

	// create uploadcare client
	uploadcareClient, err := uploadcare.New(logger, cfg.Uploadcare)
	if err != nil {
		return fmt.Errorf("creating uploadcare client: %w", err)
	}

	// create storage
	storage := repo.New(db)

	// run background tasks
	taskProvider := taskprocessor.New(logger, storage, igdbClient)
	scheduler := gocron.NewScheduler(time.UTC)
	_, err = scheduler.Cron(cfg.Scheduler.FetchIGDBGames).Do(taskProvider.StartFetchIGDBGames)
	if err != nil {
		logger.Error("run task", zap.String("task", taskprocessor.FetchIGDBGamesTaskName), zap.Error(err))
	}
	scheduler.StartAsync()

	// start debug service
	go func() {
		logger.Info("Debug service started", zap.String("address", cfg.Web.DebugAddress))
		err = http.ListenAndServe(cfg.Web.DebugAddress, nil)
		logger.Error("Debug service stopped", zap.Error(err))
	}()

	h, err := handler.Service(logger, db, auth, storage, igdbClient, uploadcareClient, cfg.Web, cfg.Zipkin)
	if err != nil {
		return fmt.Errorf("creating service handler: %w", err)
	}

	// start API service
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      h,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("API service started", zap.String("address", api.Addr))
		serverErrors <- api.ListenAndServe()
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

		err = api.Shutdown(ctx)
		if err != nil {
			logger.Error("Shutdown did not complete", zap.Duration("timeout", timeout), zap.Error(err))
			err = api.Close()
		}

		if err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}

	return nil
}
