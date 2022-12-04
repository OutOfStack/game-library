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

	"github.com/OutOfStack/game-library/internal/app/game-library-api/handler"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	conf "github.com/OutOfStack/game-library/internal/pkg/config"
	"github.com/OutOfStack/game-library/internal/pkg/database"
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
	defer logger.Sync()

	var cfg appconf.Cfg
	if err := conf.Load(".", "app", "env", &cfg); err != nil {
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
	defer db.Close()

	// create auth module
	a, err := auth.New(logger, cfg.Auth.SigningAlgorithm, cfg.Auth.VerifyTokenAPIURL)
	if err != nil {
		return fmt.Errorf("creating Auth: %w", err)
	}

	// create IGDB client
	igdbClient, err := igdb.New(logger, cfg.IGDB)
	if err != nil {
		return fmt.Errorf("creating IGDB client: %w", err)
	}

	h, err := handler.Service(logger, db, a, igdbClient, cfg.Web, cfg.Zipkin)
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

	// start debug service
	go func() {
		logger.Info("Debug service started", zap.String("address", cfg.Web.DebugAddress))
		err := http.ListenAndServe(cfg.Web.DebugAddress, nil)
		logger.Error("Debug service stopped", zap.Error(err))
	}()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("API service started", zap.String("address", api.Addr))
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("listening and serving: %w", err)
	case <-shutdown:
		logger.Info("Start shutdown")
		timeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
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
