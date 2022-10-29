package schema

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // migration file is being read here
	"github.com/jmoiron/sqlx"
)

// Migrate applies migrations
// if up is true applies all migrations otherwise rollbacks last migration
func Migrate(db *sqlx.DB, up bool) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "games", driver)
	if err != nil {
		return err
	}
	if up {
		return m.Up()
	}
	return m.Steps(-1)
}
