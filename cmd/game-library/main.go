package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/civil"
	"github.com/OutOfStack/game-library/internal/schema"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db, true); err != nil {
			log.Fatalf("applying migrations %v", err)
		}
		log.Print("migration complete")
		return
	case "rollback":
		if err := schema.Migrate(db, false); err != nil {
			log.Fatalf("rollback last migration %v", err)
		}
		log.Print("migration rollback complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatalf("applying seeds %v", err)
		}
		log.Print("Seed data inserted")
		return
	}

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

func openDB() (*sqlx.DB, error) {
	query := url.Values{}
	query.Set("sslmode", "disable")
	query.Set("timezone", "utc")

	conn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "games",
		RawQuery: query.Encode(),
	}
	log.Println(conn.String())

	return sqlx.Open("postgres", conn.String())
}

// Game represents game
type Game struct {
	Name        string     `json:"name"`
	Developer   string     `json:"developer"`
	ReleaseDate civil.Date `json:"releaseDate"`
	Genre       []string   `json:"genre"`
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
	w.Header().Set("content-type", "application/json;charset=utf-8")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error writing", err)
	}
}
