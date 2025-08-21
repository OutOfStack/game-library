package api_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_GetGames_Success() {
	page, pageSize, name, orderBy, genre, developer, publisher :=
		uint32(td.Uint8()+1), uint32(td.Uint8()+1), td.String(), "name", td.String(), td.String(), td.String()
	genreID, developerID, publisherID, platformID, count := td.Int31(), td.Int31(), td.Int31(), td.Int31(), uint64(td.Uint32())

	games := []model.Game{{
		ID:            td.Int31(),
		Name:          name,
		DevelopersIDs: []int32{developerID},
		PublishersIDs: []int32{publisherID},
		GenresIDs:     []int32{genreID},
		PlatformsIDs:  []int32{platformID},
	}}
	genresMap := map[int32]model.Genre{genreID: {
		ID:   genreID,
		Name: genre,
	}}
	companiesMap := map[int32]model.Company{
		developerID: {
			ID:   developerID,
			Name: developer,
		}, publisherID: {
			ID:   publisherID,
			Name: publisher,
		}}
	platformsMap := map[int32]model.Platform{platformID: {
		ID: platformID,
	}}

	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/games/?page=%d&pageSize=%d&name=%s&orderBy=%s&genre=%d&developer=%d&publisher=%d",
			page, pageSize, name, orderBy, genreID, developerID, publisherID),
		nil)

	s.gameFacadeMock.EXPECT().GetGames(mock.Any(), page, pageSize, model.GamesFilter{
		Name:        name,
		DeveloperID: developerID,
		PublisherID: publisherID,
		GenreID:     genreID,
		OrderBy:     repo.OrderGamesByName,
	}).Return(games, count, nil)
	s.gameFacadeMock.EXPECT().GetGenresMap(mock.Any()).Return(genresMap, nil)
	s.gameFacadeMock.EXPECT().GetPlatformsMap(mock.Any()).Return(platformsMap, nil)
	s.gameFacadeMock.EXPECT().GetCompaniesMap(mock.Any()).Return(companiesMap, nil)

	s.provider.GetGames(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(fmt.Sprintf(
		`{"games":[{"id":%d,"name":"%s","developers":[{"id":%d,"name":"%s"}],"publishers":[{"id":%d,"name":"%s"}],"releaseDate":"","genres":[{"id":%d,"name":"%s"}],"rating":0,"platforms":[{"id":%d,"name":"","abbreviation":""}],"screenshots":null,"websites":null}],"count":%d}`,
		games[0].ID, name, developerID, developer, publisherID, publisher, genreID, genre, platformID, count),
		s.httpResponse.Body.String())
}

func (s *TestSuite) Test_GetGames_InvalidFilter() {
	req := httptest.NewRequest(http.MethodGet, "/games/", nil)

	s.provider.GetGames(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetGames_Error() {
	page, pageSize := uint32(td.Uint8()+1), uint32(td.Uint8()+1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games/?page=%d&pageSize=%d", page, pageSize), nil)

	s.gameFacadeMock.EXPECT().GetGames(mock.Any(), page, pageSize, mock.Any()).Return(nil, uint64(0), errors.New("new error"))

	s.provider.GetGames(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
