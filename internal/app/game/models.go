package game

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/lib/pq"
)

// Date is date type
type Date civil.Date

// Scan casts time.Time date to Date type
func (t *Date) Scan(v interface{}) error {
	date := civil.DateOf(v.(time.Time))
	*t = Date(date)
	return nil
}

// Game represents game
type Game struct {
	ID          uint32         `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Developer   string         `db:"developer" json:"developer"`
	ReleaseDate Date           `db:"releasedate" json:"releaseDate"`
	Genre       pq.StringArray `db:"genre" json:"genre"`
}
