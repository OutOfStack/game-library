package repo

import (
	"database/sql"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents database game model
type Game struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Developer   string         `db:"developer"`
	Publisher   string         `db:"publisher"`
	ReleaseDate types.Date     `db:"release_date"`
	Price       float32        `db:"price"`
	Genre       pq.StringArray `db:"genre"`
	LogoURL     sql.NullString `db:"logo_url"`
}

// GameExt represents extended game model
type GameExt struct {
	Game
	CurrentPrice float32 `db:"current_price"`
	Rating       float32 `db:"rating"`
}

// CreateGame represents data for creating game
type CreateGame struct {
	Name        string
	Developer   string
	Publisher   string
	ReleaseDate string
	Price       float32
	Genre       []string
	LogoURL     string
}

// UpdateGame represents data for updating game
type UpdateGame struct {
	ID          int64
	Name        string
	Developer   string
	Publisher   string
	ReleaseDate string
	Price       float32
	Genre       []string
	LogoURL     string
}

// Sale represents database sale model
type Sale struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	BeginDate types.Date `db:"begin_date"`
	EndDate   types.Date `db:"end_date"`
}

// CreateSale represents data for creating sale
type CreateSale struct {
	Name      string
	BeginDate string
	EndDate   string
}

// GameSale represents GameSale model for reading from db
type GameSale struct {
	SaleID          int64  `db:"sale_id"`
	GameID          int64  `db:"game_id"`
	Sale            string `db:"sale"`
	BeginDate       string `db:"begin_date"`
	EndDate         string `db:"end_date"`
	DiscountPercent uint8  `db:"discount_percent"`
}

// CreateGameSale represents data for adding game on sale
type CreateGameSale struct {
	SaleID          int64
	GameID          int64
	DiscountPercent uint8
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
	GameID int64
}

// UserRating represents user rating entity
type UserRating struct {
	GameID int64 `db:"game_id"`
	Rating uint8 `db:"rating"`
}
