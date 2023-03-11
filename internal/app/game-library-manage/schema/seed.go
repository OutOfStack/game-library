package schema

import (
	"github.com/OutOfStack/game-library/scripts"
	"github.com/jmoiron/sqlx"
)

// Seed seeds database
func Seed(db *sqlx.DB) error {
	q := scripts.SeedSql

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(q)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
