package schema

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // register pgx5 driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // read migration file
)

const (
	migrationsSrc = "file://./scripts/migrations"
)

// Migrate applies migrations
// if up is true applies all migrations otherwise rollbacks last migration
func Migrate(dsn string, up bool) error {
	m, err := PrepareMigrations(dsn, migrationsSrc)
	if err != nil {
		return fmt.Errorf("connect to db: %v", err)
	}
	defer m.Close()

	var mErr error
	if up {
		mErr = m.Up()
	} else {
		mErr = m.Steps(-1)
	}
	if mErr != nil {
		if errors.Is(mErr, migrate.ErrNoChange) {
			log.Print(mErr)
			return nil
		}
		return mErr
	}

	return nil
}

// PrepareMigrations returns migrate instance for database migrations.
// Close() should be called after use
func PrepareMigrations(dsn, migrationsSrc string) (*migrate.Migrate, error) {
	pgxSpecificDSN := strings.Replace(dsn, "postgres://", "pgx5://", 1)

	m, err := migrate.New(migrationsSrc, pgxSpecificDSN)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %v", err)
	}
	return m, nil
}
