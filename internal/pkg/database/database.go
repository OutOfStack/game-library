package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register postgres driver
)

// Open opens connection with database
func Open() (*sqlx.DB, error) {
	query := url.Values{}
	query.Set("sslmode", "disable")
	query.Set("timezone", "utc")

	conn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "games",
		RawQuery: query.Encode(),
	}

	return sqlx.Open("postgres", conn.String())
}
