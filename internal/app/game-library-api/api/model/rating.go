package model

import (
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
)

// CreateRatingRequest - create rating request
type CreateRatingRequest struct {
	Rating uint8 `json:"rating" validate:"gte=0,lte=5"` // 0 - remove rating
}

// Validate validates CreateRatingRequest
func (r *CreateRatingRequest) Validate() (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Rating > 5 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "rating",
			Error: ErrInvalidRatingMsg,
		})
	}

	return len(validationErrors) == 0, validationErrors
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

// Validate validates GetUserRatingsRequest
func (r *GetUserRatingsRequest) Validate() (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if !validatePositive(r.GameIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "gameIds",
			Error: ErrNonPositiveValuesMsg,
		})
	}

	return len(validationErrors) == 0, validationErrors
}
