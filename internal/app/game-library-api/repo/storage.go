package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("")
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

// Storage provides required dependencies for repository
type Storage struct {
	db *sqlx.DB
}

// New creates new Storage
func New(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}
