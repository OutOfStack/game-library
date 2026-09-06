package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		implicit bool
		level    zapcore.Level
		message  string
	}{
		{name: "success", status: http.StatusCreated, level: zap.DebugLevel, message: "HTTP Request"},
		{name: "implicit success", status: http.StatusOK, implicit: true, level: zap.DebugLevel, message: "HTTP Request"},
		{name: "redirect", status: http.StatusFound, level: zap.DebugLevel, message: "HTTP Request"},
		{name: "bad request", status: http.StatusBadRequest, level: zap.DebugLevel, message: "HTTP Request"},
		{name: "unauthorized", status: http.StatusUnauthorized, level: zap.WarnLevel, message: "Unauthorized access attempt"},
		{name: "forbidden", status: http.StatusForbidden, level: zap.WarnLevel, message: "Unauthorized access attempt"},
		{name: "internal error", status: http.StatusInternalServerError, level: zap.ErrorLevel, message: "Server error"},
		{name: "unavailable", status: http.StatusServiceUnavailable, level: zap.ErrorLevel, message: "Server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			const body = "réponse"
			calls := 0
			handler := middleware.Logger(zap.New(core))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				assert.Equal(t, "query", r.URL.Query().Get("name"))
				w.Header().Set("X-Handler", "called")
				if !tt.implicit {
					w.WriteHeader(tt.status)
				}
				_, err := w.Write([]byte(body))
				assert.NoError(t, err)
			}))
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/games?name=query", nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)

			require.Equal(t, 1, calls)
			require.Equal(t, tt.status, response.Code)
			require.Equal(t, body, response.Body.String())
			require.Equal(t, "called", response.Header().Get("X-Handler"))
			entries := logs.All()
			require.Len(t, entries, 1)
			require.Equal(t, tt.level, entries[0].Level)
			require.Equal(t, tt.message, entries[0].Message)
			fields := entries[0].ContextMap()
			require.Equal(t, int64(tt.status), fields["status"])
			require.Equal(t, http.MethodGet, fields["method"])
			require.Equal(t, "/games", fields["path"])
			duration, ok := fields["duration"].(time.Duration)
			require.True(t, ok)
			require.GreaterOrEqual(t, duration, time.Duration(0))
			if tt.level == zap.DebugLevel {
				require.Equal(t, int64(len(body)), fields["bytes"])
			}
		})
	}
}
