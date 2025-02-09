package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetIDParam(t *testing.T) {
	tests := []struct {
		name               string
		url                string
		expectedID         int32
		expectedStatusCode int
	}{
		{
			name:               "Valid ID",
			url:                "/resource/123",
			expectedID:         123,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid ID (non-numeric)",
			url:                "/resource/abc",
			expectedID:         0,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid ID (negative)",
			url:                "/resource/-1",
			expectedID:         0,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid ID (zero)",
			url:                "/resource/0",
			expectedID:         0,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Missing ID",
			url:                "/resource/",
			expectedID:         0,
			expectedStatusCode: http.StatusNotFound, // Default Chi behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/resource/{id}", func(w http.ResponseWriter, r *http.Request) {
				id, err := web.GetIDParam(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				w.Write([]byte(string(id))) //nolint
			})

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedStatusCode == http.StatusOK {
				assert.Equal(t, string(tt.expectedID), rr.Body.String())
			}
		})
	}
}
