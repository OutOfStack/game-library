package api_test

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetPlatforms_Success() {
	s.gameFacadeMock.EXPECT().GetPlatforms(mock.Any()).Return([]model.Platform{
		{ID: 1, Name: "PlayStation 5", Abbreviation: "PS5"},
		{ID: 2, Name: "Linux", Abbreviation: "Linux"},
	}, nil)

	s.provider.GetPlatforms(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(`[
			{"id": 1, "name": "PlayStation 5", "abbreviation": "PS5"},
			{"id": 2, "name": "Linux", "abbreviation": "Linux"}
		]`, s.httpResponse.Body.String())
}

func (s *TestSuite) TestGetPlatforms_Error() {
	s.gameFacadeMock.EXPECT().GetPlatforms(mock.Any()).Return(nil, errors.New("new error"))

	s.provider.GetPlatforms(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
