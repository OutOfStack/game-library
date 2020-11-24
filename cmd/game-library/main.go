package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer log.Println("Completed")

	api := http.Server{
		Addr:         ":8000",
		Handler:      http.HandlerFunc(Echo),
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

func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Method ", r.Method, "\nPath ", r.URL.Path)
}
