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

func (s *TestSuite) Test_GetUserRatings_Success() {
	gameID, rating, authToken, userID, role := td.Int31(), uint8(td.Intn(5)+1), td.String(), td.String(), td.String()

	requestData := api.GetUserRatingsRequest{
		GameIDs: []int32{gameID},
	}
	ratings := map[int32]uint8{
		gameID: rating,
	}

	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/user/ratings/", bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetUserRatings(mock.Any(), userID).Return(ratings, nil)

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetUserRatings)))
	r := chi.NewRouter()
	r.Post("/user/ratings/", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(fmt.Sprintf(`{"%d":%d}`, gameID, rating), s.httpResponse.Body.String())
}

func (s *TestSuite) Test_GetUserRatings_MissingClaims() {
	requestData := api.GetUserRatingsRequest{}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/user/ratings/", bytes.NewReader(requestBody))

	s.provider.GetUserRatings(s.httpResponse, req)

	fmt.Println(s.httpResponse.Body.String())
	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetUserRatings_FacadeError() {
	authToken, userID, role := td.String(), td.String(), td.String()

	requestData := api.GetUserRatingsRequest{}

	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/user/ratings/", bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetUserRatings(mock.Any(), userID).Return(nil, errors.New("new error"))

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetUserRatings)))
	r := chi.NewRouter()
	r.Post("/user/ratings/", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
