package web_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/OutOfStack/game-library/internal/web"
	"github.com/stretchr/testify/require"
)

func TestErrorConstructors(t *testing.T) {
	cause := errors.New("invalid game")
	fields := []web.FieldError{{Field: "name", Error: "required"}}
	tests := []struct {
		name    string
		err     error
		status  int
		message string
		fields  []web.FieldError
		cause   error
	}{
		{
			name: "existing error", err: web.NewError(cause, http.StatusBadRequest),
			status: http.StatusBadRequest, message: "invalid game", cause: cause,
		},
		{
			name: "message", err: web.NewErrorFromMessage("missing game", http.StatusNotFound),
			status: http.StatusNotFound, message: "missing game",
		},
		{
			name: "status", err: web.NewErrorFromStatusCode(http.StatusServiceUnavailable),
			status: http.StatusServiceUnavailable, message: "Service Unavailable",
		},
		{
			name: "fields", err: web.NewErrorWithFields(cause, http.StatusUnprocessableEntity, fields),
			status: http.StatusUnprocessableEntity, message: "invalid game", fields: fields, cause: cause,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := errors.AsType[*web.Error](tt.err)
			require.True(t, ok)
			require.Equal(t, tt.status, got.StatusCode)
			require.EqualError(t, got.Err, tt.message)
			require.Equal(t, tt.fields, got.Fields)
			if tt.cause != nil {
				require.ErrorIs(t, got.Err, tt.cause)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name   string
		fields []web.FieldError
		want   string
	}{
		{name: "nil fields", want: "invalid game"},
		{name: "empty fields", fields: []web.FieldError{}, want: "invalid game"},
		{
			name: "multiple fields", fields: []web.FieldError{{Field: "name", Error: "required"}, {Field: "rating", Error: "invalid"}},
			want: "invalid game - fields: [{name required} {rating invalid}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &web.Error{Err: errors.New("invalid game"), Fields: tt.fields}
			require.EqualError(t, err, tt.want)
		})
	}
}
