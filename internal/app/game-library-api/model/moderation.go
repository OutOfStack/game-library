package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// ModerationStatus represents moderation status of a game
type ModerationStatus string

const (
	// ModerationStatusPending represents a game that needs moderation
	ModerationStatusPending ModerationStatus = "pending"
	// ModerationStatusInProgress represents a game that is being moderated at the moment
	ModerationStatusInProgress ModerationStatus = "in_progress"
	// ModerationStatusReady represents a game that is ready
	ModerationStatusReady ModerationStatus = "ready"
	// ModerationStatusDeclined represents a game that is declined and requires fixing
	ModerationStatusDeclined ModerationStatus = "declined"
)

// Moderation represents stored moderation record
type Moderation struct {
	ID        int32          `db:"id"`
	GameID    int32          `db:"game_id"`
	Status    string         `db:"status"`
	Details   string         `db:"details"`
	Error     sql.NullString `db:"error"`
	GameData  ModerationData `db:"game_data"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

// ModerationData represents game data for moderation
type ModerationData struct {
	Name        string   `json:"name"`
	Developers  []string `json:"developers"`
	Publisher   string   `json:"publisher"`
	ReleaseDate string   `json:"releaseDate"`
	Genres      []string `json:"genres"`
	LogoURL     string   `json:"logoUrl"`
	Summary     string   `json:"summary"`
	Slug        string   `json:"slug"`
	Screenshots []string `json:"screenshots"`
	Websites    []string `json:"websites"`
}

// GameModerationData aggregates a game with its moderation history
type GameModerationData struct {
	Game        Game
	Moderations []Moderation
}

// CreateModeration represents data required to create moderation record
type CreateModeration struct {
	GameID   int32
	GameData ModerationData
	Status   ModerationStatus
}

// NewCreateModeration creates new CreateModeration
func NewCreateModeration(gameID int32, gameData ModerationData) CreateModeration {
	return CreateModeration{
		GameID:   gameID,
		GameData: gameData,
		Status:   ModerationStatusPending,
	}
}

// UpdateModerationResult represents data to update moderation result
type UpdateModerationResult struct {
	ResultStatus ModerationStatus
	Details      string
	Error        sql.NullString
}

// Value implements driver.Valuer for ModerationData
func (md ModerationData) Value() (driver.Value, error) {
	b, err := json.Marshal(md)
	if err != nil {
		return nil, fmt.Errorf("marshal ModerationData: %w", err)
	}
	return string(b), nil
}

// Scan implements sql.Scanner for ModerationData
func (md *ModerationData) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), md)
	case []byte:
		return json.Unmarshal(v, md)
	default:
		return fmt.Errorf("scan ModerationData: unsupported type %T", src)
	}
}
