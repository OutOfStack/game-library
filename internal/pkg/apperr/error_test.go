package apperr_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError(t *testing.T) {
	id := td.Int64()
	entity := td.String()

	err := apperr.NewNotFoundError(entity, id)

	assert.Equal(t, entity, err.Entity)
	assert.Equal(t, id, err.ID)
	assert.Equal(t, apperr.NotFound, err.StatusCode())
	assert.Equal(t, http.StatusNotFound, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("%s with id %d not found", entity, id), err.Error())
}

func TestNewInvalidError(t *testing.T) {
	id := td.Int64()
	entity := td.String()
	msg := td.String()

	err := apperr.NewInvalidError(entity, id, msg)

	assert.Equal(t, entity, err.Entity)
	assert.Equal(t, id, err.ID)
	assert.Equal(t, apperr.Invalid, err.StatusCode())
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("invalid %s with id %d: %s", entity, id, msg), err.Error())
}

func TestNewForbiddenError(t *testing.T) {
	id := td.Int64()
	entity := td.String()

	err := apperr.NewForbiddenError(entity, id)

	assert.Equal(t, entity, err.Entity)
	assert.Equal(t, id, err.ID)
	assert.Equal(t, apperr.Forbidden, err.StatusCode())
	assert.Equal(t, http.StatusForbidden, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("forbidden access to %s with id %v", entity, id), err.Error())
}

func TestError_Methods(t *testing.T) {
	id := td.String()
	entity := td.String()
	msg := td.String()

	err := apperr.Error[string]{
		Entity:   entity,
		ID:       id,
		StatCode: apperr.Invalid,
		Msg:      msg,
	}

	assert.Equal(t, apperr.Invalid, err.StatusCode())
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("invalid %s with id %s: %s", entity, id, msg), err.Error())
}

func TestIsAppError_True(t *testing.T) {
	id := td.Int32()
	entity := td.String()

	err := apperr.NewNotFoundError(entity, id)

	// Direct check
	appErr, ok := apperr.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperr.NotFound, appErr.StatusCode())
	assert.Equal(t, http.StatusNotFound, appErr.HTTPStatusCode())

	// Wrapped error check
	wrappedErr := fmt.Errorf("wrapped error: %w", err)
	appErr, ok = apperr.IsAppError(wrappedErr)
	assert.True(t, ok)
	assert.Equal(t, apperr.NotFound, appErr.StatusCode())
	assert.Equal(t, http.StatusNotFound, appErr.HTTPStatusCode())
}

func TestIsAppError_False(t *testing.T) {
	err := errors.New("some other error")

	// Direct check
	appErr, ok := apperr.IsAppError(err)
	assert.False(t, ok)
	assert.Nil(t, appErr)
}

func TestIsStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode apperr.StatusCode
		expected   bool
	}{
		{
			name:       "AppError with expected status code",
			err:        apperr.NewNotFoundError(td.String(), td.Int64()),
			statusCode: apperr.NotFound,
			expected:   true,
		},
		{
			name:       "AppError with different status code",
			err:        apperr.NewInvalidError(td.String(), td.String(), td.String()),
			statusCode: apperr.NotFound,
			expected:   false,
		},
		{
			name:       "Non-AppError",
			err:        errors.New(td.String()),
			statusCode: apperr.NotFound,
			expected:   false,
		},
		{
			name:       "Wrapped AppError with expected status code",
			err:        fmt.Errorf("wrapped error: %w", apperr.NewInvalidError(td.String(), td.Int32(), td.String())),
			statusCode: apperr.Invalid,
			expected:   true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := apperr.IsStatusCode(tt.err, tt.statusCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}
