package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/authapi"
	"github.com/OutOfStack/game-library/internal/middleware"
	mwmock "github.com/OutOfStack/game-library/internal/middleware/mocks"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type AuthTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	log            *zap.Logger
	authClientMock *mwmock.MockAuthClient
}

func (s *AuthTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.log = zap.NewNop()
	s.authClientMock = mwmock.NewMockAuthClient(s.ctrl)
}

func (s *AuthTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestAuthTestSuite_Run(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

// Authenticate middleware tests
func (s *AuthTestSuite) Test_Authenticate_Success() {
	token := td.String()
	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *AuthTestSuite) Test_Authenticate_NoAuthorizationHeader() {
	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
	s.Contains(rr.Body.String(), "no Authorization header found")
}

func (s *AuthTestSuite) Test_Authenticate_NoBearerToken() {
	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc123")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
	s.Contains(rr.Body.String(), "no Bearer token found")
}

func (s *AuthTestSuite) Test_Authenticate_InvalidToken() {
	token := td.String()
	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(errors.New("invalid token"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
}

func (s *AuthTestSuite) Test_Authenticate_EmptyBearerToken() {
	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer ")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
	s.Contains(rr.Body.String(), "no Bearer token found")
}

func (s *AuthTestSuite) Test_Authenticate_VerifyAPIUnavailable() {
	token := td.String()
	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(authapi.ErrVerifyAPIUnavailable)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusBadGateway, rr.Code)
}

// Authorize middleware tests
func (s *AuthTestSuite) Test_Authorize_Success() {
	token := td.String()
	requiredRole := "admin"
	claims := &auth.Claims{
		Name:     td.String(),
		UserRole: requiredRole,
	}

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, requiredRole)
	handler := authenticator(authorizer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(nil)
	s.authClientMock.EXPECT().ParseToken(token).Return(claims, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *AuthTestSuite) Test_Authorize_MissingToken() {
	requiredRole := "admin"
	authorizer := middleware.Authorize(s.log, s.authClientMock, requiredRole)
	handler := authorizer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

func (s *AuthTestSuite) Test_Authorize_InvalidRole() {
	token := td.String()
	requiredRole := "admin"
	userRole := "user"
	claims := &auth.Claims{
		Name:     td.String(),
		UserRole: userRole,
	}

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, requiredRole)
	handler := authenticator(authorizer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(nil)
	s.authClientMock.EXPECT().ParseToken(token).Return(claims, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusForbidden, rr.Code)
	s.Contains(rr.Body.String(), "access denied")
}

func (s *AuthTestSuite) Test_Authorize_ParseTokenError() {
	token := td.String()
	requiredRole := "admin"

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, requiredRole)
	handler := authenticator(authorizer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(nil)
	s.authClientMock.EXPECT().ParseToken(token).Return(nil, errors.New("invalid token format"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusInternalServerError, rr.Code)
}

// GetClaims tests
func (s *AuthTestSuite) Test_GetClaims_Success() {
	token := td.String()
	requiredRole := "admin"
	claims := &auth.Claims{
		Name:     td.String(),
		UserRole: requiredRole,
	}

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, requiredRole)

	handler := authenticator(authorizer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedClaims, err := middleware.GetClaims(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if retrievedClaims.UserRole != requiredRole {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	})))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(nil)
	s.authClientMock.EXPECT().ParseToken(token).Return(claims, nil)

	req := httptest.NewRequest(http.MethodPost, "/games", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}

func (s *AuthTestSuite) Test_GetClaims_NotFound() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := middleware.GetClaims(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusUnauthorized, rr.Code)
}

// AuthenticationFlow Integration tests
func (s *AuthTestSuite) Test_AuthenticationFlow_MultipleRoles() {
	token1 := td.String()
	token2 := td.String()
	publisherRole := "publisher"
	moderatorRole := "moderator"

	publisherClaims := &auth.Claims{
		Name:     td.String(),
		UserRole: publisherRole,
	}
	moderatorClaims := &auth.Claims{
		Name:     td.String(),
		UserRole: moderatorRole,
	}

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	publisherAuthorizer := middleware.Authorize(s.log, s.authClientMock, publisherRole)

	handler := authenticator(publisherAuthorizer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token1).Return(nil)
	s.authClientMock.EXPECT().ParseToken(token1).Return(publisherClaims, nil)

	req1 := httptest.NewRequest(http.MethodPost, "/games", nil)
	req1.Header.Set("Authorization", "Bearer "+token1)
	rr1 := httptest.NewRecorder()

	handler.ServeHTTP(rr1, req1)

	s.Equal(http.StatusOK, rr1.Code)

	s.authClientMock.EXPECT().Verify(gomock.Any(), token2).Return(nil)
	s.authClientMock.EXPECT().ParseToken(token2).Return(moderatorClaims, nil)

	req2 := httptest.NewRequest(http.MethodPost, "/games", nil)
	req2.Header.Set("Authorization", "Bearer "+token2)
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr2, req2)

	s.Equal(http.StatusForbidden, rr2.Code)
}

func (s *AuthTestSuite) Test_Authenticate_BearerTokenPassedToHandler() {
	token := td.String()

	handler := middleware.Authenticate(s.log, s.authClientMock)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedClaims, err := middleware.GetClaims(r.Context())
		if err == nil && retrievedClaims != nil {
			_ = retrievedClaims.Name
		}
		w.WriteHeader(http.StatusOK)
	}))

	s.authClientMock.EXPECT().Verify(gomock.Any(), token).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
}
