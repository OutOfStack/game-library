package model_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/stretchr/testify/require"
)

func TestCreateRatingRequestValidation(t *testing.T) {
	t.Run("Valid rating", func(t *testing.T) {
		request := model.CreateRatingRequest{
			Rating: 4,
		}

		valid, errors := request.Validate()
		require.True(t, valid, "Expected valid request")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Zero rating is valid", func(t *testing.T) {
		request := model.CreateRatingRequest{
			Rating: 0,
		}

		valid, errors := request.Validate()
		require.True(t, valid, "Expected valid request")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Rating too high", func(t *testing.T) {
		request := model.CreateRatingRequest{
			Rating: 6, // Above max of 5
		}

		valid, errors := request.Validate()
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 1, "Expected 1 validation error")
		require.Equal(t, "rating", errors[0].Field)
		require.Equal(t, model.ErrInvalidRatingMsg, errors[0].Error)
	})
}

func TestGetUserRatingsRequestValidation(t *testing.T) {
	t.Run("Valid GameIDs", func(t *testing.T) {
		request := model.GetUserRatingsRequest{
			GameIDs: []int32{1, 2, 3},
		}

		valid, errors := request.Validate()
		require.True(t, valid, "Expected valid request")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Empty GameIDs is valid", func(t *testing.T) {
		request := model.GetUserRatingsRequest{
			GameIDs: []int32{},
		}

		valid, errors := request.Validate()
		require.True(t, valid, "Expected valid request")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Non-positive GameIDs", func(t *testing.T) {
		request := model.GetUserRatingsRequest{
			GameIDs: []int32{1, 0, 3}, // Contains zero
		}

		valid, errors := request.Validate()
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 1, "Expected 1 validation error")
		require.Equal(t, "gameIds", errors[0].Field)
		require.Equal(t, model.ErrNonPositiveValuesMsg, errors[0].Error)
	})
}
