package game

import (
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/lib/pq"
)

// Game represents game
type Game struct {
	ID          int            `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Developer   string         `db:"developer" json:"developer"`
	ReleaseDate types.Date     `db:"release_date" json:"releaseDate"`
	Genre       pq.StringArray `db:"genre" json:"genre"`
}

// PostModel represents game model we receive from user
// TODO: change releaseDate type to time.TIme (requires custom JSON marsaller and unmarhsller)
type PostModel struct {
	Name        string         `json:"name"`
	Developer   string         `json:"developer"`
	ReleaseDate string         `json:"releaseDate"`
	Genre       pq.StringArray `json:"genre"`
}

// Sale respresents information about game being on sale
type Sale struct {
	ID              int        `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	GameID          int        `db:"game_id" json:"gameId"`
	BeginDate       types.Date `db:"begin_date" json:"beginDate"`
	EndDate         types.Date `db:"end_date" json:"endDate"`
	DiscountPercent uint8      `db:"discount_percent" json:"discountPercent"`
}
