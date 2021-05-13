package game

import (
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents database game model
type Game struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Developer   string         `db:"developer"`
	ReleaseDate types.Date     `db:"release_date"`
	Price       float32        `db:"price"`
	Genre       pq.StringArray `db:"genre"`
}

// GetGame represents game get model
type GetGame struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Developer   string   `json:"developer"`
	ReleaseDate string   `json:"releaseDate"`
	Price       float32  `json:"price"`
	Genre       []string `json:"genre"`
}

// CreateGame represents game data we receive from user
type CreateGame struct {
	Name        string   `json:"name" validate:"required"`
	Developer   string   `json:"developer" validate:"required"`
	ReleaseDate string   `json:"releaseDate" validate:"date"`
	Price       float32  `json:"price" validate:"gte=0,lt=10000"`
	Genre       []string `json:"genre"`
}

// UpdateGame represents model for updating information about game.
// All fields are optional
type UpdateGame struct {
	Name        *string   `json:"name"`
	Developer   *string   `json:"developer" validate:"omitempty"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty,date"`
	Price       *float32  `json:"price" validate:"omitempty,gte=0,lt=10000"`
	Genre       *[]string `json:"genre" validate:"omitempty"`
}

// Sale respresents database sale model
type Sale struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	BeginDate types.Date `db:"begin_date"`
	EndDate   types.Date `db:"end_date"`
}

// GetSale respresents sale get model
type GetSale struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

// CreateSale represents sale data we receive from user
type CreateSale struct {
	Name      string `json:"name"`
	BeginDate string `json:"beginDate" validate:"required,date"`
	EndDate   string `json:"endDate" validate:"required,date"`
}

// GameSale represents GameSale model for reading from db
type GameSale struct {
	SaleID          int64 `db:"sale_id"`
	GameID          int64 `db:"game_id"`
	Sale            string
	BeginDate       string `db:"begin_date"`
	EndDate         string `db:"end_date"`
	DiscountPercent uint8  `db:"discount_percent"`
}

// CreateGameSale represents data about game being on sale
type CreateGameSale struct {
	SaleID          int64 `json:"saleId"`
	DiscountPercent uint8 `json:"discountPercent" validate:"gt=0,lte=100"`
}

// GetGameSale represents game sale get model
type GetGameSale struct {
	GameID          int64  `json:"gameId"`
	SaleID          int64  `json:"saleId"`
	Sale            string `json:"sale"`
	DiscountPercent uint8  `json:"discountPercent"`
	BeginDate       string `json:"beginDate"`
	EndDate         string `json:"endDate"`
}
