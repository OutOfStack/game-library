package schema

import (
	"context"

	"github.com/OutOfStack/game-library/scripts"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Seed seeds database
func Seed(db *pgxpool.Pool) error {
	q := scripts.SeedSQL

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, q); err != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			return rErr
		}
		return err
	}

	return tx.Commit(ctx)
}
