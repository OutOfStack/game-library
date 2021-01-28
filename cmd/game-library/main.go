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
	"github.com/lib/pq"
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

	svc := GameService{db: db}

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

// Date is date type
type Date civil.Date

// Scan casts time.Time date to Date type
func (t *Date) Scan(v interface{}) error {
	date := civil.DateOf(v.(time.Time))
	*t = Date(date)
	return nil
}

// Game represents game
type Game struct {
	ID          uint32         `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Developer   string         `db:"developer" json:"developer"`
	ReleaseDate Date           `db:"releasedate" json:"releaseDate"`
	Genre       pq.StringArray `db:"genre" json:"genre"`
}

// GameService has handler methods for dealing with games
type GameService struct {
	db *sqlx.DB
}

// List returns all games
func (g *GameService) List(w http.ResponseWriter, r *http.Request) {
	list := []Game{}

	const q = `select id, name, developer, releasedate, genre from games`

	if err := g.db.Select(&list, q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error querying db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json;charset=utf-8")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error writing", err)
	}
}
