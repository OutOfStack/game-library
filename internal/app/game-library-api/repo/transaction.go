package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// BeginTx starts transaction
func (s *Storage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}
