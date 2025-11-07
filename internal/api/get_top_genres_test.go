package api_test

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/model"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestTopGetGenres_Success() {
	s.gameFacadeMock.EXPECT().GetTopGenres(mock.Any(), int64(10)).Return([]model.Genre{
		{ID: 1, Name: "RPG"},
		{ID: 2, Name: "Quest"},
	}, nil)

	s.provider.GetTopGenres(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(`[
			{"id": 1, "name": "RPG"},
			{"id": 2, "name": "Quest"}
		]`, s.httpResponse.Body.String())
}

func (s *TestSuite) TestTopGetGenres_Error() {
	s.gameFacadeMock.EXPECT().GetTopGenres(mock.Any(), int64(10)).Return(nil, errors.New("new error"))

	s.provider.GetTopGenres(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
