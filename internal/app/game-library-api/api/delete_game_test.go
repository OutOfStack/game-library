package api_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/go-chi/chi/v5"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_DeleteGame_Success() {
	gameID, authToken, publisher, role := td.Int31(), td.String(), td.String(), td.String()

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/games/%d", gameID), nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: publisher, UserRole: role}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().DeleteGame(mock.Any(), gameID, publisher).Return(nil)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.DeleteGame)))
	r := chi.NewRouter()
	r.Delete("/games/{id}", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusNoContent, s.httpResponse.Code)
}

func (s *TestSuite) Test_DeleteGame_InvalidID() {
	req := httptest.NewRequest(http.MethodDelete, "/games/-100", nil)

	r := chi.NewRouter()
	r.Delete("/games/{id}", s.provider.DeleteGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
	s.JSONEq(`{"error":"invalid id"}`, s.httpResponse.Body.String())
}

func (s *TestSuite) Test_DeleteGame_MissingClaims() {
	gameID := td.Int31()

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/games/%d", gameID), nil)

	handler := http.HandlerFunc(s.provider.DeleteGame)
	r := chi.NewRouter()
	r.Delete("/games/{id}", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_DeleteGame_FacadeError() {
	gameID, authToken, publisher, role := td.Int31(), td.String(), td.String(), td.String()

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/games/%d", gameID), nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: publisher, UserRole: role}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().DeleteGame(mock.Any(), gameID, publisher).Return(errors.New("new error"))

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.DeleteGame)))
	r := chi.NewRouter()
	r.Delete("/games/{id}", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
