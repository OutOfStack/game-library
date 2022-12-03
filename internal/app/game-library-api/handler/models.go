package handler

import (
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
)

// GameResp represents game response
type GameResp struct {
	ID          int32    `json:"id"`
	Name        string   `json:"name"`
	Developer   string   `json:"developer"`
	Publisher   string   `json:"publisher"`
	ReleaseDate string   `json:"releaseDate"`
	Genre       []string `json:"genre"`
	LogoURL     string   `json:"logoUrl,omitempty"`
	Rating      float32  `json:"rating"`
}

// CreateGameReq represents game data we receive from user
type CreateGameReq struct {
	Name        string   `json:"name" validate:"required"`
	Developer   string   `json:"developer" validate:"required"`
	ReleaseDate string   `json:"releaseDate" validate:"date"`
	Genre       []string `json:"genre"`
	LogoURL     string   `json:"logoUrl"`
}

// UpdateGameReq represents model for updating information about game.\
// All fields are optional
type UpdateGameReq struct {
	Name        *string   `json:"name"`
	Developer   *string   `json:"developer" validate:"omitempty"`
	Publisher   *string   `json:"publisher" validate:"omitempty"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty,date"`
	Genre       *[]string `json:"genre" validate:"omitempty"`
	LogoURL     *string   `json:"logoUrl"`
}

// CreateRatingReq represents rating data we receive from user
type CreateRatingReq struct {
	Rating uint8 `json:"rating" validate:"gte=1,lte=5"`
}

// RatingResp represents response to rating request
type RatingResp struct {
	GameID int32 `json:"gameId"`
	Rating uint8 `json:"rating"`
}

// UserRatingsReq represents get user ratings request
type UserRatingsReq struct {
	GameIDs []int32 `json:"gameIds"`
}

// IDResp represents response with id
type IDResp struct {
	ID int32 `json:"id"`
}

func mapToCreateRating(crr *CreateRatingReq, gameID int32, userID string) repo.CreateRating {
	return repo.CreateRating{
		Rating: crr.Rating,
		UserID: userID,
		GameID: gameID,
	}
}

func mapCreateRatingToResp(cr repo.CreateRating) *RatingResp {
	return &RatingResp{
		GameID: cr.GameID,
		Rating: cr.Rating,
	}
}

func mapToCreateGame(cgr *CreateGameReq) repo.CreateGame {
	return repo.CreateGame{
		Name:        cgr.Name,
		Developer:   cgr.Developer,
		ReleaseDate: cgr.ReleaseDate,
		Genre:       cgr.Genre,
		LogoURL:     cgr.LogoURL,
	}
}

func mapToUpdateGame(g repo.Game, ugr UpdateGameReq) repo.UpdateGame {
	var logoURL string
	if g.LogoURL.Valid {
		logoURL = g.LogoURL.String
	}

	update := repo.UpdateGame{
		Name:        g.Name,
		Developer:   g.Developer,
		Publisher:   g.Publisher,
		ReleaseDate: g.ReleaseDate.String(),
		LogoURL:     logoURL,
		Genre:       g.Genre,
	}

	if ugr.Name != nil {
		update.Name = *ugr.Name
	}
	if ugr.Developer != nil {
		update.Developer = *ugr.Developer
	}
	if ugr.Publisher != nil {
		update.Publisher = *ugr.Publisher
	}
	if ugr.ReleaseDate != nil {
		update.ReleaseDate = *ugr.ReleaseDate
	}
	if ugr.Genre != nil {
		update.Genre = *ugr.Genre
	}
	if ugr.LogoURL != nil && *ugr.LogoURL != "" {
		update.LogoURL = *ugr.LogoURL
	}

	return update
}

func mapGameToResp(g repo.Game) GameResp {
	return GameResp{
		ID:          g.ID,
		Name:        g.Name,
		Developer:   g.Developer,
		Publisher:   g.Publisher,
		ReleaseDate: g.ReleaseDate.String(),
		Genre:       []string(g.Genre),
		LogoURL:     g.LogoURL.String,
		Rating:      float32(g.Rating.Float64),
	}
}
