package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/handler"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/pkg/util"
)

func main() {

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
			ReadTimeout     time.Duration `mapstructure:"APP_READTIMEOUT"`
			WriteTimeout    time.Duration `mapstructure:"APP_WRITETIMEOUT"`
			ShutdownTimeout time.Duration `mapstructure:"APP_SHUTDOWNTIMEOUT"`
		} `mapstructure:",squash"`
	}

	cfg := config{}
	if err := util.LoadConfig(".", "app", "env", &cfg); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	//fmt.Printf("%+v\n", cfg)

	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		RequireSSL: cfg.DB.RequireSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	svc := handler.Game{DB: db}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(svc.List),
		ReadTimeout:  cfg.Web.ReadTimeout * time.Second,
		WriteTimeout: cfg.Web.WriteTimeout * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Error listening and serving: %v", err)
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
			log.Fatalf("could not stop server: %v", err)
		}
	}
}
