package model

// CreateRating represents data for rating a game
type CreateRating struct {
	Rating uint8
	UserID string
	GameID int32
}

// RemoveRating represents data for removing game rating
type RemoveRating struct {
	UserID string
	GameID int32
}

// UserRating represents user rating entity
type UserRating struct {
	GameID int32  `db:"game_id"`
	UserID string `db:"user_id"`
	Rating uint8  `db:"rating"`
}
