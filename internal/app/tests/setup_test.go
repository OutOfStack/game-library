package tests_test

import (
	"fmt"
	"log"
	"os"
	"testing"

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
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
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
		log.Fatalf("Could not start docker container: %v", err)
	}
	log.Println("Docker container started")

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	counter := 1
	if err := pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:%s@localhost:%s/games?sslmode=disable", PgPwd, HostPort))
		if err != nil {
			log.Printf("Attempt %d connecting to database: %v", counter, err)
			counter++
			return err
		}
		err = db.Ping()
		if err != nil {
			log.Printf("Attempt %d pinging database: %v", counter, err)
			counter++
		}
		return err
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not connect to database: %s", err)
	}
	log.Println("Database connection established")

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	log.Println("Container deleted")

	os.Exit(code)
}

func TestDbConnection(t *testing.T) {
	t.Logf("Testing postgres connection with docker\n")

	err := db.Ping()
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	t.Logf("Succeeded\n")
}
