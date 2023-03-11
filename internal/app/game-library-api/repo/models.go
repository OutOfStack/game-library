package repo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents database game model
type Game struct {
	ID          int32          `db:"id"`
	Name        string         `db:"name"`
	Developers  pq.Int32Array  `db:"developers"`
	Publishers  pq.Int32Array  `db:"publishers"`
	ReleaseDate types.Date     `db:"release_date"`
	Genres      pq.Int32Array  `db:"genres"`
	LogoURL     string         `db:"logo_url"`
	Rating      float64        `db:"rating"`
	Summary     string         `db:"summary"`
	Slug        string         `db:"slug"`
	Platforms   pq.Int32Array  `db:"platforms"`
	Screenshots pq.StringArray `db:"screenshots"`
	Websites    pq.StringArray `db:"websites"`
	IGDBRating  float64        `db:"igdb_rating"`
	IGDBID      int64          `db:"igdb_id"`
	Weight      float64        `db:"weight"` // Readonly field
}

// CreateGame represents data for creating game
type CreateGame struct {
	Name        string
	Developers  []int32
	Publishers  []int32
	ReleaseDate string
	Genres      []int32
	LogoURL     string
	Summary     string
	Slug        string
	Platforms   []int32
	Screenshots []string
	Websites    []string
	IGDBRating  float64
	IGDBID      int64
}

// UpdateGame represents data for updating game
type UpdateGame struct {
	Name        string
	Developers  []int32
	Publishers  []int32
	ReleaseDate string
	Genres      []int32
	LogoURL     string
	Summary     string
	Slug        string
	Platforms   []int32
	Screenshots []string
	Websites    []string
	IGDBRating  float64
	IGDBID      int64
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
	Settings TaskSettings `db:"settings"`
}

// TaskSettings task settings value
type TaskSettings []byte

// Value implements driver.Valuer interface
func (ts TaskSettings) Value() (driver.Value, error) {
	return string(ts), nil
}

// Scan implements sql.Scanner interface
func (ts *TaskSettings) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		*ts = []byte(v)
	case []byte:
		*ts = v
	default:
		return fmt.Errorf("scan TaskSettings: unsupported type %T", src)
	}

	return nil
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
	ID     int32         `db:"id"`
	Name   string        `db:"name"`
	IGDBID sql.NullInt64 `db:"igdb_id"`
}
