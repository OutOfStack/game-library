package api_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/go-chi/chi/v5"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_GetGame_Success() {
	id, name := td.Int31(), td.String()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games/%d", id), nil)

	s.gameFacadeMock.EXPECT().GetGameByID(mock.Any(), id).Return(model.Game{ID: id, Name: name}, nil)
	s.gameFacadeMock.EXPECT().GetGenresMap(mock.Any()).Return(nil, nil)
	s.gameFacadeMock.EXPECT().GetPlatformsMap(mock.Any()).Return(nil, nil)
	s.gameFacadeMock.EXPECT().GetCompaniesMap(mock.Any()).Return(nil, nil)

	r := chi.NewRouter()
	r.Get("/games/{id}", s.provider.GetGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(fmt.Sprintf(
		`{"id":%d,"name":"%s","developers":null,"publishers":null,"releaseDate":"","genres":null,"rating":0,"platforms":null,"screenshots":null,"websites":null}`, id, name),
		s.httpResponse.Body.String())
}

func (s *TestSuite) Test_GetGame_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/games/-100", nil)

	r := chi.NewRouter()
	r.Get("/games/{id}", s.provider.GetGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetGame_NotFound() {
	id := td.Int31()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games/%d", id), nil)

	s.gameFacadeMock.EXPECT().GetGameByID(mock.Any(), id).Return(model.Game{}, apperr.NewNotFoundError("game", id))

	r := chi.NewRouter()
	r.Get("/games/{id}", s.provider.GetGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusNotFound, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetGame_Error() {
	id := td.Int31()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games/%d", id), nil)

	s.gameFacadeMock.EXPECT().GetGameByID(mock.Any(), id).Return(model.Game{}, errors.New("new error"))

	r := chi.NewRouter()
	r.Get("/games/{id}", s.provider.GetGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
