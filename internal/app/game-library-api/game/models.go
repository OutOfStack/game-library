package game

import (
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents game
type Game struct {
	ID          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Developer   string         `db:"developer" json:"developer"`
	ReleaseDate types.Date     `db:"release_date" json:"releaseDate"`
	Price       float32        `db:"price" json:"price"`
	Genre       pq.StringArray `db:"genre" json:"genre"`
}

// NewGame represents game data we receive from user
// TODO: change releaseDate type to time.TIme (requires custom JSON marsaller and unmarhsller)
type NewGame struct {
	Name        string         `json:"name" validate:"required"`
	Developer   string         `json:"developer" validate:"required"`
	ReleaseDate string         `json:"releaseDate"`
	Price       float32        `json:"price" validate:"gte=0"`
	Genre       pq.StringArray `json:"genre"`
}

// UpdtaeGame represents model for updating information about game.
// All fields are optional
type UpdateGame struct {
	Name        *string         `json:"name"`
	Developer   *string         `json:"developer" validate:"omitempty"`
	ReleaseDate *string         `json:"releaseDate" validate:"omitempty"`
	Price       *float32        `json:"price" validate:"gte=0"`
	Genre       *pq.StringArray `json:"genre" validate:"omitempty"`
}

// Sale respresents information about game being on sale
type Sale struct {
	ID              int64      `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	GameID          int64      `db:"game_id" json:"gameId"`
	BeginDate       types.Date `db:"begin_date" json:"beginDate"`
	EndDate         types.Date `db:"end_date" json:"endDate"`
	DiscountPercent uint8      `db:"discount_percent" json:"discountPercent"`
}

// NewSale represents sale data we receive from user
type NewSale struct {
	Name            string `json:"name"`
	BeginDate       string `json:"beginDate" validate:"required"`
	EndDate         string `json:"endDate" validate:"required"`
	DiscountPercent uint8  `json:"discountPercent" validate:"gt=0,lte=100"`
}
