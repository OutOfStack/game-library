package apperr

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError(t *testing.T) {
	id := td.Int64()
	entity := td.String()

	err := NewNotFoundError(entity, id)

	assert.Equal(t, entity, err.Entity)
	assert.Equal(t, id, err.ID)
	assert.Equal(t, NotFound, err.StatusCode())
	assert.Equal(t, http.StatusNotFound, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("%s with id %d not found", entity, id), err.Error())
}

func TestNewInvalidError(t *testing.T) {
	id := td.Int64()
	entity := td.String()
	msg := td.String()

	err := NewInvalidError(entity, id, msg)

	assert.Equal(t, entity, err.Entity)
	assert.Equal(t, id, err.ID)
	assert.Equal(t, Invalid, err.StatusCode())
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("invalid %s with id %d: %s", entity, id, msg), err.Error())
}

func TestError_Methods(t *testing.T) {
	id := td.String()
	entity := td.String()
	msg := td.String()

	err := Error[string]{
		Entity:   entity,
		ID:       id,
		StatCode: Invalid,
		Msg:      msg,
	}

	assert.Equal(t, Invalid, err.StatusCode())
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatusCode())
	assert.Equal(t, fmt.Sprintf("invalid %s with id %s: %s", entity, id, msg), err.Error())
}

func TestIsAppError_True(t *testing.T) {
	id := td.Int32()
	entity := td.String()

	err := NewNotFoundError(entity, id)

	// Direct check
	appErr, ok := IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, NotFound, appErr.StatusCode())
	assert.Equal(t, http.StatusNotFound, appErr.HTTPStatusCode())

	// Wrapped error check
	wrappedErr := fmt.Errorf("wrapped error: %w", err)
	appErr, ok = IsAppError(wrappedErr)
	assert.True(t, ok)
	assert.Equal(t, NotFound, appErr.StatusCode())
	assert.Equal(t, http.StatusNotFound, appErr.HTTPStatusCode())
}

func TestIsAppError_False(t *testing.T) {
	err := errors.New("some other error")

	// Direct check
	appErr, ok := IsAppError(err)
	assert.False(t, ok)
	assert.Nil(t, appErr)
}

func TestIsStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode StatusCode
		expected   bool
	}{
		{
			name:       "AppError with expected status code",
			err:        NewNotFoundError(td.String(), td.Int64()),
			statusCode: NotFound,
			expected:   true,
		},
		{
			name:       "AppError with different status code",
			err:        NewInvalidError(td.String(), td.String(), td.String()),
			statusCode: NotFound,
			expected:   false,
		},
		{
			name:       "Non-AppError",
			err:        errors.New(td.String()),
			statusCode: NotFound,
			expected:   false,
		},
		{
			name:       "Wrapped AppError with expected status code",
			err:        fmt.Errorf("wrapped error: %w", NewInvalidError(td.String(), td.Int32(), td.String())),
			statusCode: Invalid,
			expected:   true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStatusCode(tt.err, tt.statusCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}
