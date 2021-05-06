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
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/pkg/util"
)

// @title Game library API
// @version 0.1
// @description API for game library service
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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
	log := log.New(os.Stdout, "GAMES : ", log.LstdFlags)

	type config struct {
		DB struct {
			Host       string `mapstructure:"APP_HOST"`
			Name       string `mapstructure:"APP_NAME"`
			User       string `mapstructure:"APP_USER"`
			Password   string `mapstructure:"APP_PASSWORD"`
			RequireSSL bool   `mapstructure:"APP_REQUIRESSL"`
		} `mapstructure:",squash"`
		Web struct {
			Address         string        `mapstructure:"APP_ADDRESS"`
			Debug           string        `mapstructure:"DEBUG"`
			ReadTimeout     time.Duration `mapstructure:"APP_READTIMEOUT"`
			WriteTimeout    time.Duration `mapstructure:"APP_WRITETIMEOUT"`
			ShutdownTimeout time.Duration `mapstructure:"APP_SHUTDOWNTIMEOUT"`
		} `mapstructure:",squash"`
	}

	cfg := config{}
	if err := util.LoadConfig(".", "app", "env", &cfg); err != nil {
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

	// start debug service
	go func() {
		log.Printf("Debug service listening on %s", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, nil)
		log.Printf("Debug service stopped %v", err)
	}()

	// start API service
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handler.Service(log, db),
		ReadTimeout:  cfg.Web.ReadTimeout * time.Second,
		WriteTimeout: cfg.Web.WriteTimeout * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("API service listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("listening and serving: %w", err)
	case <-shutdown:
		log.Println("Start shutdown")
		timeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("Shutdown did not complete in %s : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}

	return nil
}
