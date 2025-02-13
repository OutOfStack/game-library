package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New returns database connection
func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create database connection pool: %v", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping database: %v", err)
	}

	return pool, nil
}
