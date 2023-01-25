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

// CommitTx commits transaction
func (s *Storage) CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}

// RollbackTx rollbacks transaction
func (s *Storage) RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}
