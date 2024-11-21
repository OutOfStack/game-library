package api_test

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetPlatforms_Success() {
	s.GameFacadeMock.EXPECT().GetPlatforms(gomock.Any()).Return([]model.Platform{
		{ID: 1, Name: "PlayStation 5", Abbreviation: "PS5"},
		{ID: 2, Name: "Linux", Abbreviation: "Linux"},
	}, nil)

	// Act
	s.Provider.GetPlatforms(s.HTTPResponse, s.HTTPRequest)

	// Assert
	s.Equal(http.StatusOK, s.HTTPResponse.Code)
	s.JSONEq(`[
			{"id": 1, "name": "PlayStation 5", "abbreviation": "PS5"},
			{"id": 2, "name": "Linux", "abbreviation": "Linux"}
		]`, s.HTTPResponse.Body.String())
}

func (s *TestSuite) TestGetPlatforms_Error() {
	s.GameFacadeMock.EXPECT().GetPlatforms(gomock.Any()).Return(nil, errors.New("new error"))

	// Act
	s.Provider.GetPlatforms(s.HTTPResponse, s.HTTPRequest)

	// Assert
	s.Equal(http.StatusInternalServerError, s.HTTPResponse.Code)
}
