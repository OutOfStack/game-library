package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// BeginTx starts transaction
func (s *Storage) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, &sql.TxOptions{})
}
