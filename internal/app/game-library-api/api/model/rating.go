package model

// CreateRatingRequest - create rating request
type CreateRatingRequest struct {
	Rating uint8 `json:"rating" validate:"gte=0,lte=5"` // 0 - remove rating
}

// RatingResponse - rating response
type RatingResponse struct {
	GameID int32 `json:"gameId"`
	Rating uint8 `json:"rating"`
}

// GetUserRatingsRequest - get user ratings request
type GetUserRatingsRequest struct {
	GameIDs []int32 `json:"gameIds"`
}
