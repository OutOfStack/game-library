package repo

import (
	"database/sql"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents database game model
type Game struct {
	ID          int32          `db:"id"`
	Name        string         `db:"name"`
	Developer   string         `db:"developer"` // Deprecated: use Developers instead
	Developers  pq.Int32Array  `db:"developers"`
	Publisher   string         `db:"publisher"` // Deprecated: use Publishers instead
	Publishers  pq.Int32Array  `db:"publishers"`
	ReleaseDate types.Date     `db:"release_date"`
	Genre       pq.StringArray `db:"genre"` // Deprecated: use Genres instead
	Genres      pq.Int32Array  `db:"genres"`
	LogoURL     string         `db:"logo_url"`
	Rating      float64        `db:"rating"`
	Summary     string         `db:"summary"`
	Slug        string         `db:"slug"`
	Platforms   pq.Int32Array  `db:"platforms"`
	Screenshots pq.StringArray `db:"screenshots"`
	Websites    pq.StringArray `db:"websites"`
	IGDBRating  float64        `db:"igdb_rating"`
	IGDBID      float64        `db:"igdb_id"`
}

// CreateGame represents data for creating game
type CreateGame struct {
	Name        string
	Developer   string
	Publisher   string
	ReleaseDate string
	Genre       []string
	LogoURL     string
}

// UpdateGame represents data for updating game
type UpdateGame struct {
	Name        string
	Developer   string
	Publisher   string
	ReleaseDate string
	Genre       []string
	LogoURL     string
}

// CreateRating represents data for rating a game
type CreateRating struct {
	Rating uint8
	UserID string
	GameID int32
}

// UserRating represents user rating entity
type UserRating struct {
	GameID int32  `db:"game_id"`
	UserID string `db:"user_id"`
	Rating uint8  `db:"rating"`
}

// TaskStatus represents task status type
type TaskStatus string

// Task status values
const (
	IdleTaskStatus    TaskStatus = "idle"
	RunningTaskStatus TaskStatus = "running"
	ErrorTaskStatus   TaskStatus = "error"
)

// Task represents task entity
type Task struct {
	Name     string       `db:"name"`
	Status   TaskStatus   `db:"status"`
	RunCount int64        `db:"run_count"`
	LastRun  sql.NullTime `db:"last_run"`
}

// Platform represents platform entity
type Platform struct {
	ID           int32  `db:"id"`
	Name         string `db:"name"`
	Abbreviation string `db:"abbreviation"`
	IGDBID       int64  `db:"igdb_id"`
}

// Genre represents genre entity
type Genre struct {
	ID     int32  `db:"id"`
	Name   string `db:"name"`
	IGDBID int64  `db:"igdb_id"`
}

// Company represents company entity
type Company struct {
	ID     int32  `db:"id"`
	Name   string `db:"name"`
	IGDBID int64  `db:"igdb_id"`
}
