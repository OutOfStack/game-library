package web_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/stretchr/testify/assert"
)

func TestRespond(t *testing.T) {
	tests := []struct {
		name       string
		val        interface{}
		statusCode int
		expected   string
	}{
		{
			name:       "Respond with valid JSON",
			val:        map[string]string{"key": "value"},
			statusCode: http.StatusOK,
			expected:   `{"key":"value"}`,
		},
		{
			name:       "Respond with nil value",
			val:        nil,
			statusCode: http.StatusOK,
			expected:   `{}`,
		},
		{
			name:       "Respond with No Content",
			val:        nil,
			statusCode: http.StatusNoContent,
			expected:   "",
		},
		{
			name:       "Respond with invalid value",
			val:        make(chan int), // Invalid type for JSON marshaling
			statusCode: http.StatusInternalServerError,
			expected:   `{"error":"marshal value to json:`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			web.Respond(rr, tt.val, tt.statusCode)

			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			if tt.expected != "" {
				assert.Contains(t, string(body), tt.expected)
			} else {
				assert.Empty(t, string(body))
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Internal server error with nil error",
			err:            nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal Server Error"}`,
		},
		{
			name:           "Custom error with client error status",
			err:            web.NewErrorFromMessage("bad request", http.StatusBadRequest),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"bad request"}`,
		},
		{
			name:           "Custom error with server error status",
			err:            web.NewErrorFromStatusCode(http.StatusServiceUnavailable),
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   `{"error":"Service Unavailable"}`,
		},
		{
			name:           "Error is not of type *web.Error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			web.RespondError(rr, tt.err)

			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}

func TestRespond500(t *testing.T) {
	rr := httptest.NewRecorder()
	web.Respond500(rr)

	resp := rr.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Internal Server Error"}`, string(body))
}
