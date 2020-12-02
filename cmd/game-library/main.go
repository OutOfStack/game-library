package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/civil"
)

func main() {
	defer log.Println("Completed")

	api := http.Server{
		Addr:         ":8000",
		Handler:      http.HandlerFunc(ListGames),
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

// Game represents game
type Game struct {
	Name        string
	Developer   string
	ReleaseDate civil.Date
	Genre       []string
}

// ListGames returns all games
func ListGames(w http.ResponseWriter, r *http.Request) {
	list := []Game{
		{
			Name:        "Red Dead Redemption 2",
			Developer:   "Rockstar Games",
			Genre:       []string{"Action", "Western", "Adventure"},
			ReleaseDate: civil.Date{Year: 2019, Month: 12, Day: 5},
		},
		{
			Name:        "Ori and the Will of the Wisps",
			Developer:   "Moon Studios GmbH",
			Genre:       []string{"Action", "Platformer"},
			ReleaseDate: civil.Date{Year: 2020, Month: 3, Day: 11},
		},
		{
			Name:        "The Wolf Among Us",
			Developer:   "Telltale",
			Genre:       []string{"Adventure", "Episodic", "Detective"},
			ReleaseDate: civil.Date{Year: 2013, Month: 10, Day: 11},
		},
	}
	data, err := json.Marshal(list)
	if err != nil {
		log.Println("Error marshalling", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error writing", err)
	}
}
