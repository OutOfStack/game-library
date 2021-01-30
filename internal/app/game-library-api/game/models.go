package game

import (
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents game
type Game struct {
	ID          uint32         `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Developer   string         `db:"developer" json:"developer"`
	ReleaseDate types.Date     `db:"releasedate" json:"releaseDate"`
	Genre       pq.StringArray `db:"genre" json:"genre"`
}
