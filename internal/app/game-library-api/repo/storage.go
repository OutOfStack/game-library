package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("db")
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

// Storage provides required dependencies for repository
type Storage struct {
	db *pgxpool.Pool
}

// New creates new Storage
func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		db: db,
	}
}
