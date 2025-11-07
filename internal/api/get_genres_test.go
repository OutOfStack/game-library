package api_test

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/model"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetGenres_Success() {
	s.gameFacadeMock.EXPECT().GetGenres(mock.Any()).Return([]model.Genre{
		{ID: 1, Name: "RPG"},
		{ID: 2, Name: "Quest"},
	}, nil)

	s.provider.GetGenres(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(`[
			{"id": 1, "name": "RPG"},
			{"id": 2, "name": "Quest"}
		]`, s.httpResponse.Body.String())
}

func (s *TestSuite) TestGetGenres_Error() {
	s.gameFacadeMock.EXPECT().GetGenres(mock.Any()).Return(nil, errors.New("new error"))

	s.provider.GetGenres(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
