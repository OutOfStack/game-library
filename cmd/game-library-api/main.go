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
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	svc := handler.Game{DB: db}

	api := http.Server{
		Addr:         ":8000",
		Handler:      http.HandlerFunc(svc.List),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
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
		const timeout = 5 * time.Second
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
