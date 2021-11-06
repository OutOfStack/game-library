package game

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
	LogoUrl     sql.NullString `db:"logo_url"`
}

// GameInfo represents extended info about game
type GameInfo struct {
	Game
	CurrentPrice float32 `db:"current_price"`
	Rating       float32 `db:"rating"`
}

// GameResp represents game response
type GameResp struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Developer   string   `json:"developer"`
	Publisher   string   `json:"publisher"`
	ReleaseDate string   `json:"releaseDate"`
	Price       float32  `json:"price"`
	Genre       []string `json:"genre"`
	LogoUrl     string   `json:"logoUrl,omitempty"`
}

// GameInfoResp represents extended game info response
type GameInfoResp struct {
	GameResp
	CurrentPrice float32 `json:"currentPrice"`
	Rating       float32 `json:"rating"`
}

// CreateGame represents game data we receive from user
type CreateGameReq struct {
	Name        string   `json:"name" validate:"required"`
	Developer   string   `json:"developer" validate:"required"`
	Publisher   string   `json:"-"`
	ReleaseDate string   `json:"releaseDate" validate:"date"`
	Price       float32  `json:"price" validate:"gte=0,lt=10000"`
	Genre       []string `json:"genre"`
	LogoUrl     string   `json:"logoUrl"`
}

// UpdateGame represents model for updating information about game.
// All fields are optional
type UpdateGameReq struct {
	Name        *string   `json:"name"`
	Developer   *string   `json:"developer" validate:"omitempty"`
	Publisher   *string   `json:"publisher" validate:"omitempty"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty,date"`
	Price       *float32  `json:"price" validate:"omitempty,gte=0,lt=10000"`
	Genre       *[]string `json:"genre" validate:"omitempty"`
	LogoUrl     *string   `json:"logoUrl"`
}

// Sale represents database sale model
type Sale struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	BeginDate types.Date `db:"begin_date"`
	EndDate   types.Date `db:"end_date"`
}

// SaleResp represents sale response
type SaleResp struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

// CreateSale represents sale data we receive from user
type CreateSaleReq struct {
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
type CreateGameSaleReq struct {
	SaleID          int64 `json:"saleId"`
	DiscountPercent uint8 `json:"discountPercent" validate:"gt=0,lte=100"`
}

// GameSaleResp represents game sale response
type GameSaleResp struct {
	GameID          int64  `json:"gameId"`
	SaleID          int64  `json:"saleId"`
	Sale            string `json:"sale"`
	DiscountPercent uint8  `json:"discountPercent"`
	BeginDate       string `json:"beginDate"`
	EndDate         string `json:"endDate"`
}

// Rating represents database rating model
type Rating struct {
	GameID int64  `db:"game_id"`
	UserID string `db:"user_id"`
	Rating uint8  `db:"rating"`
}

// CreateRating represents rating data we receive from user
type CreateRatingReq struct {
	Rating uint8 `json:"rating" validate:"gte=1,lte=4"`
}

// RatingResp represents response to rating request
type RatingResp struct {
	GameID int64 `json:"gameId"`
	Rating uint8 `json:"rating"`
}

// UserRating represents user rating entity
type UserRating struct {
	GameID int64 `db:"game_id"`
	Rating uint8 `db:"rating"`
}

// UserRatings represents get user ratings request
type UserRatingsReq struct {
	GameIDs []int64 `json:"gameIds"`
}
