package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_GetUserGames_Success() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()

	games := []model.Game{
		{
			ID:          td.Int31(),
			Name:        td.String(),
			ReleaseDate: types.DateOf(td.Date()),
			Summary:     td.String(),
			LogoURL:     td.String(),
		},
	}

	expectedResponse := []api.GameResponse{
		{
			ID:          games[0].ID,
			Name:        games[0].Name,
			ReleaseDate: games[0].ReleaseDate.String(),
			Summary:     games[0].Summary,
			LogoURL:     games[0].LogoURL,
			Developers:  nil,
			Publishers:  nil,
			Genres:      nil,
			Platforms:   nil,
			Screenshots: nil,
			Websites:    nil,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/user/games", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetPublisherGames(mock.Any(), publisherName).Return(games, nil)
	s.gameFacadeMock.EXPECT().GetGenresMap(mock.Any()).Return(map[int32]model.Genre{}, nil)
	s.gameFacadeMock.EXPECT().GetCompaniesMap(mock.Any()).Return(map[int32]model.Company{}, nil)
	s.gameFacadeMock.EXPECT().GetPlatformsMap(mock.Any()).Return(map[int32]model.Platform{}, nil)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetUserGames)))
	r := chi.NewRouter()
	r.Get("/user/games", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)

	var response []api.GameResponse
	err := json.Unmarshal(s.httpResponse.Body.Bytes(), &response)
	s.NoError(err)
	s.Equal(expectedResponse, response)
}

func (s *TestSuite) Test_GetUserGames_MissingClaims() {
	req := httptest.NewRequest(http.MethodGet, "/user/games", nil)

	s.provider.GetUserGames(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetUserGames_FacadeError() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()

	req := httptest.NewRequest(http.MethodGet, "/user/games", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetPublisherGames(mock.Any(), publisherName).Return(nil, errors.New("facade error"))

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetUserGames)))
	r := chi.NewRouter()
	r.Get("/user/games", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_GetUserGames_MapToGameResponseError() {
	authToken, userID, role := td.String(), td.String(), td.String()
	publisherName := td.String()

	games := []model.Game{
		{
			ID:   td.Int31(),
			Name: td.String(),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/user/games", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, UserRole: role, Name: publisherName}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().GetPublisherGames(mock.Any(), publisherName).Return(games, nil)
	s.gameFacadeMock.EXPECT().GetGenresMap(mock.Any()).Return(nil, errors.New("genres error"))

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.GetUserGames)))
	r := chi.NewRouter()
	r.Get("/user/games", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
