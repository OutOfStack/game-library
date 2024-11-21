package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestExtractToken(t *testing.T) {
	s1, s2 := td.String(), td.String()
	tests := []struct {
		name       string
		authHeader string
		expected   string
	}{
		{"Valid Bearer Token", "Bearer " + s1, s1},
		{"Invalid Format", "Bearer", ""},
		{"Empty Header", "", ""},
		{"Non-Bearer Token", "Basic " + s2, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := auth.ExtractToken(tt.authHeader)
			require.Equal(t, tt.expected, token)
		})
	}
}

func TestClient_Verify_Success(t *testing.T) {
	token := td.String()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody auth.VerifyToken
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, token, reqBody.Token)

		respBody := auth.VerifyTokenResp{Valid: true}
		_ = json.NewEncoder(w).Encode(respBody)
	}))
	defer server.Close()

	logger := zap.NewNop()
	client, _ := auth.New(logger, jwt.SigningMethodHS256.Alg(), server.URL, nil)

	err := client.Verify(context.Background(), token)
	require.NoError(t, err)
}

func TestClient_Verify_InvalidToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		respBody := auth.VerifyTokenResp{Valid: false}
		_ = json.NewEncoder(w).Encode(respBody)
	}))
	defer server.Close()

	logger := zap.NewNop()
	client, _ := auth.New(logger, jwt.SigningMethodHS256.Alg(), server.URL, nil)

	err := client.Verify(context.Background(), td.String())
	require.Error(t, err)
	require.EqualError(t, err, "invalid token")
}

func TestClient_Verify_APIUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(10 * time.Millisecond) // more than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpClient := &http.Client{
		Timeout: 1 * time.Millisecond,
	}

	logger := zap.NewNop()
	client, _ := auth.New(logger, jwt.SigningMethodHS256.Alg(), server.URL, httpClient)

	err := client.Verify(context.Background(), td.String())
	require.Error(t, err)
	require.Equal(t, auth.ErrVerifyAPIUnavailable, err)
}

func TestClient_ParseToken(t *testing.T) {
	logger := zap.NewNop()
	client, _ := auth.New(logger, jwt.SigningMethodHS256.Alg(), "", nil)

	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcl9yb2xlIjoibW9kZXJhdG9yIiwidXNlcm5hbWUiOiJqb2huZG9lIiwibmFtZSI6IkpvaG4gRG9lIn0.abc" //nolint
	claims, err := client.ParseToken(tokenStr)

	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, "1234567890", claims.UserID())
	require.Equal(t, auth.RoleModerator, claims.UserRole)
	require.Equal(t, "johndoe", claims.Username)
	require.Equal(t, "John Doe", claims.Name)
}

func TestClient_ParseToken_Invalid(t *testing.T) {
	logger := zap.NewNop()
	client, _ := auth.New(logger, jwt.SigningMethodHS256.Alg(), "", nil)

	tokenStr := "invalid-token"
	claims, err := client.ParseToken(tokenStr)

	require.Error(t, err)
	require.Nil(t, claims)
	require.Contains(t, err.Error(), "parsing token")
}
