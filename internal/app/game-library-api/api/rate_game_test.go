package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_RateGame_Success() {
	gameID, rating, authToken, userID, role := td.Int31(), uint8(td.Intn(6)), td.String(), td.String(), td.String()

	requestData := api.CreateRatingRequest{
		Rating: rating,
	}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%d/rate", gameID), bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().RateGame(mock.Any(), gameID, userID, rating).Return(nil)

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.RateGame)))
	r := chi.NewRouter()
	r.Post("/{id}/rate", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(fmt.Sprintf(`{"gameId":%d,"rating":%d}`, gameID, rating), s.httpResponse.Body.String())
}

func (s *TestSuite) Test_RateGame_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/-100/rate", nil)

	r := chi.NewRouter()
	r.Post("/{id}/rate", s.provider.RateGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
	s.JSONEq(`{"error":"invalid id"}`, s.httpResponse.Body.String())
}

func (s *TestSuite) Test_RateGame_MissingClaims() {
	gameID := td.Int31()

	requestData := api.CreateRatingRequest{}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%d/rate", gameID), bytes.NewReader(requestBody))

	handler := http.HandlerFunc(s.provider.RateGame)
	r := chi.NewRouter()
	r.Post("/{id}/rate", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_RateGame_FacadeError() {
	gameID, rating, authToken, userID, role := td.Int31(), uint8(td.Intn(6)), td.String(), td.String(), td.String()

	requestData := api.CreateRatingRequest{
		Rating: rating,
	}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%d/rate", gameID), bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().RateGame(mock.Any(), gameID, userID, rating).Return(errors.New("new error"))

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.RateGame)))
	r := chi.NewRouter()
	r.Post("/{id}/rate", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
