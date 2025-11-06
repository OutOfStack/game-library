package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type txKey struct{}

// Querier defines the interface for executing queries.
// It's implemented by *pgxpool.Pool and pgx.Tx.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// BeginTx starts transaction
func (s *Storage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}

// WithTx stores transaction in context
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromContext retrieves transaction from context
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// RunWithTx runs a function in a transaction
func (s *Storage) RunWithTx(ctx context.Context, f func(context.Context) error) error {
	// check if we're already in a transaction
	if _, ok := TxFromContext(ctx); ok {
		// just exec func
		return f(ctx)
	}

	// begin new tx
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}

	txCtx := WithTx(ctx, tx)
	err = f(txCtx)
	if err != nil {
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			s.log.Error("rollback transaction", zap.Error(txErr))
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

// querier returns appropriate Querier (transaction if available, otherwise connection pool)
func (s *Storage) querier(ctx context.Context) Querier {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return s.db
}
