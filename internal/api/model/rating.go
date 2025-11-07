package model

import (
	"github.com/OutOfStack/game-library/internal/api/validation"
	"github.com/OutOfStack/game-library/internal/web"
)

// CreateRatingRequest - create rating request
type CreateRatingRequest struct {
	Rating uint8 `json:"rating" validate:"gte=0,lte=5"` // 0 - remove rating
}

// ValidateWith validates CreateRatingRequest
func (r *CreateRatingRequest) ValidateWith(v *validation.Validator) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Rating > 5 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "rating",
			Error: v.ErrInvalidRatingMsg(),
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

// ValidateWith validates GetUserRatingsRequest
func (r *GetUserRatingsRequest) ValidateWith(v *validation.Validator) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if !v.ValidatePositive(r.GameIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "gameIds",
			Error: v.ErrNonPositiveValuesMsg(),
		})
	}

	return len(validationErrors) == 0, validationErrors
}
