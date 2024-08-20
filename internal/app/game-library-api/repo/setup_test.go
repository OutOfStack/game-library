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
	DatabaseName  = "games"
	DatabasePort  = "5439"
	DatabasePwd   = "password"
	MigrationsSrc = "file://../../../../scripts/migrations"
	pg            = "postgres"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Repo tests: Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: pg,
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=" + DatabasePwd,
			"POSTGRES_DB=" + DatabaseName,
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {
				docker.PortBinding{
					HostIP:   "0.0.0.0",
					HostPort: DatabasePort,
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
		db, err = sqlx.Open(pg, fmt.Sprintf("postgres://postgres:%s@localhost:%s/%s?sslmode=disable", DatabasePwd, DatabasePort, DatabaseName))
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
	t.Helper()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		t.Fatalf("error on creating db driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(MigrationsSrc, "games", driver)
	if err != nil {
		t.Fatalf("error on connecting to db: %v", err)
	}

	if err = m.Up(); err != nil {
		t.Fatalf("error on applying migrations: %v", err)
	}
	return repo.New(db)
}

func teardown(t *testing.T) {
	t.Helper()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		t.Fatalf("error on creating db driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(MigrationsSrc, "games", driver)
	if err != nil {
		t.Fatalf("error on connecting to db: %v", err)
	}

	if err = m.Down(); err != nil {
		t.Fatalf("error on migration rollback: %v", err)
	}
}
