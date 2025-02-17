package auth_test

import (
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/authapi"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestExtractToken() {
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
		s.Run(tt.name, func() {
			token := auth.ExtractToken(tt.authHeader)
			s.Equal(tt.expected, token)
		})
	}
}

func (s *TestSuite) TestClient_Verify_Success() {
	token := td.String()
	resp := authapi.VerifyTokenResp{Valid: true}

	s.authAPIClient.EXPECT().VerifyToken(gomock.Any(), token).Return(resp, nil)

	err := s.auth.Verify(s.ctx, token)

	s.NoError(err)
}

func (s *TestSuite) TestClient_Verify_InvalidToken() {
	token := td.String()
	resp := authapi.VerifyTokenResp{Valid: false}

	s.authAPIClient.EXPECT().VerifyToken(gomock.Any(), token).Return(resp, nil)

	err := s.auth.Verify(s.ctx, token)

	s.Require().Error(err)
	s.EqualError(err, "invalid token")
}

func (s *TestSuite) TestClient_Verify_APIUnavailable() {
	token := td.String()

	s.authAPIClient.EXPECT().VerifyToken(gomock.Any(), token).Return(authapi.VerifyTokenResp{}, authapi.ErrVerifyAPIUnavailable)

	err := s.auth.Verify(s.ctx, token)

	s.ErrorIs(err, authapi.ErrVerifyAPIUnavailable)
}

func (s *TestSuite) TestClient_ParseToken() {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcl9yb2xlIjoibW9kZXJhdG9yIiwidXNlcm5hbWUiOiJqb2huZG9lIiwibmFtZSI6IkpvaG4gRG9lIn0.abc" //nolint

	claims, err := s.auth.ParseToken(tokenStr)

	s.Require().NoError(err)
	s.NotNil(claims)
	s.Equal("1234567890", claims.UserID())
	s.Equal(auth.RoleModerator, claims.UserRole)
	s.Equal("johndoe", claims.Username)
	s.Equal("John Doe", claims.Name)
}

func (s *TestSuite) TestClient_ParseToken_Invalid() {
	tokenStr := "invalid-token"

	claims, err := s.auth.ParseToken(tokenStr)

	s.Require().Error(err)
	s.Nil(claims)
	s.Contains(err.Error(), "parsing token")
}
