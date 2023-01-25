package repo

import (
	"database/sql"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents database game model
type Game struct {
	ID          int32           `db:"id"`
	Name        string          `db:"name"`
	Developer   string          `db:"developer"`
	Publisher   string          `db:"publisher"`
	ReleaseDate types.Date      `db:"release_date"`
	Genre       pq.StringArray  `db:"genre"`
	LogoURL     sql.NullString  `db:"logo_url"`
	Rating      sql.NullFloat64 `db:"rating"`
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

// Rating represents rating entity
type Rating struct {
	GameID int64  `db:"game_id"`
	UserID string `db:"user_id"`
	Rating uint8  `db:"rating"`
}

// CreateRating represents data for rating a game
type CreateRating struct {
	Rating uint8
	UserID string
	GameID int32
}

// UserRating represents user rating entity
type UserRating struct {
	GameID int32 `db:"game_id"`
	Rating uint8 `db:"rating"`
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
