package api_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_GetGameModerations_Success() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()
	gameID := td.Int31()
	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour)

	moderations := []model.Moderation{
		{
			ID:        td.Int31(),
			Status:    "approved",
			Details:   "Game approved",
			CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
			UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		},
	}

	expectedResponse := []api.ModerationItem{
		{
			ID:        moderations[0].ID,
			Status:    moderations[0].Status,
			Details:   moderations[0].Details,
			CreatedAt: createdAt.Format(time.RFC3339),
			UpdatedAt: updatedAt.Format(time.RFC3339),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/games/"+strconv.Itoa(int(gameID))+"/moderations", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetGameModerations(mock.Any(), gameID, publisherName).Return(moderations, nil)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetGameModerations)))
	r := chi.NewRouter()
	r.Get("/games/{id}/moderations", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)

	var response []api.ModerationItem
	err := json.Unmarshal(s.httpResponse.Body.Bytes(), &response)
	s.NoError(err)
	s.Equal(expectedResponse, response)
}

func (s *TestSuite) Test_GetGameModerations_MissingClaims() {
	gameID := td.Int31()
	req := httptest.NewRequest(http.MethodGet, "/games/"+strconv.Itoa(int(gameID))+"/moderations", nil)

	s.provider.GetGameModerations(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetGameModerations_InvalidID() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()

	req := httptest.NewRequest(http.MethodGet, "/games/invalid/moderations", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetGameModerations)))
	r := chi.NewRouter()
	r.Get("/games/{id}/moderations", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetGameModerations_AppError() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()
	gameID := td.Int31()

	req := httptest.NewRequest(http.MethodGet, "/games/"+strconv.Itoa(int(gameID))+"/moderations", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	appErr := apperr.NewNotFoundError("game", gameID)
	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetGameModerations(mock.Any(), gameID, publisherName).Return(nil, appErr)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetGameModerations)))
	r := chi.NewRouter()
	r.Get("/games/{id}/moderations", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusNotFound, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetGameModerations_FacadeError() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()
	gameID := td.Int31()

	req := httptest.NewRequest(http.MethodGet, "/games/"+strconv.Itoa(int(gameID))+"/moderations", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetGameModerations(mock.Any(), gameID, publisherName).Return(nil, errors.New("facade error"))

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetGameModerations)))
	r := chi.NewRouter()
	r.Get("/games/{id}/moderations", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
