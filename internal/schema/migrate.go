package schema

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// migration file is being read here
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// Migrate applies migrations
// if up is true migrates everything othrwise rollbacks last migration
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
		err = m.Up()
	} else {
		err = m.Steps(-1)
	}

	if err != nil {
		return err
	}
	return nil
}
