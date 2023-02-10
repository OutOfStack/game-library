package repo_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	PgPwd    string = "temp_pwd"
	HostPort string = "5439"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Repo tests: Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", PgPwd),
			"POSTGRES_DB=games",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {
				docker.PortBinding{
					HostIP:   "0.0.0.0",
					HostPort: HostPort,
				},
			},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Repo tests: Could not start docker container: %v", err)
	}
	log.Println("Repo tests: Docker container started")

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	counter := 1
	err = pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:%s@localhost:%s/games?sslmode=disable", PgPwd, HostPort))
		if err != nil {
			log.Printf("Repo tests: Attempt %d connecting to database: %v", counter, err)
			counter++
			return err
		}
		err = db.Ping()
		if err != nil {
			log.Printf("Repo tests: Attempt %d pinging database: %v", counter, err)
			counter++
		}
		return err
	})
	if err != nil {
		if pErr := pool.Purge(resource); pErr != nil {
			log.Fatalf("Repo tests: Could not purge resource: %s", pErr)
		}
		log.Fatalf("Repo tests: Could not connect to database: %v", err)
	}

	log.Println("Repo tests: Database connection established")

	// runs tests in current package
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Repo tests: Could not purge resource: %s", err)
	}
	log.Println("Repo tests: Docker container deleted")

	os.Exit(code)
}

func setup(t *testing.T) *repo.Storage {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		t.Fatalf("error on creating db driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://../../../../scripts/migrations", "games", driver)
	if err != nil {
		t.Fatalf("error on connecting to db: %v", err)
	}

	if err = m.Up(); err != nil {
		t.Fatalf("error on applying migrations: %v", err)
	}
	return repo.New(db)
}

func teardown(t *testing.T) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		t.Fatalf("error on creating db driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://../../../../scripts/migrations", "games", driver)
	if err != nil {
		t.Fatalf("error on connecting to db: %v", err)
	}

	if err = m.Down(); err != nil {
		t.Fatalf("error on migration rollback: %v", err)
	}
}
