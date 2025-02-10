package repo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-manage/schema"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // registers pgx5 driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	DatabaseName  = "games"
	DatabasePort  = "5439"
	DatabaseUser  = "games-user"
	DatabasePwd   = "games-password"
	MigrationsSrc = "file://../../../../scripts/migrations"
	pg            = "postgres"
)

var dsn = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", DatabaseUser, DatabasePwd, DatabasePort, DatabaseName)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Repo tests: Could not connect to docker: %s", err)
	}
	pool.MaxWait = 30 * time.Second

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: pg,
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_USER=" + DatabaseUser,
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
	ctx := context.Background()
	err = pool.Retry(func() error {
		db, err = database.New(ctx, dsn)
		if err != nil {
			log.Printf("Repo tests: Attempt %d connecting to database", counter)
			counter++
			return err
		}
		return nil
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

	db.Close()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Repo tests: Could not purge resource: %s", err)
	}
	log.Println("Repo tests: Docker container deleted")

	os.Exit(code)
}

func setup(t *testing.T) *repo.Storage {
	t.Helper()

	m, err := schema.PrepareMigrations(dsn, MigrationsSrc)
	if err != nil {
		t.Fatalf("error preparing migrations: %v", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		t.Fatalf("error on applying migrations: %v", err)
	}
	return repo.New(db)
}

func teardown(t *testing.T) {
	t.Helper()

	m, err := schema.PrepareMigrations(dsn, MigrationsSrc)
	if err != nil {
		t.Fatalf("error preparing migrations: %v", err)
	}
	defer m.Close()

	if err = m.Down(); err != nil {
		t.Fatalf("error on migration rollback: %v", err)
	}
}
