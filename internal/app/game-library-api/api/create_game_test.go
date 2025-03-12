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
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_CreateGame_Success() {
	role, gameID, authToken := td.String(), td.Int32(), td.String()

	requestData := api.CreateGameRequest{
		Name:         td.String(),
		Developer:    td.String(),
		ReleaseDate:  td.Date().Format("2006-01-02"),
		GenresIDs:    []int32{td.Int31()},
		LogoURL:      s.getImageURL(),
		Summary:      td.String(),
		PlatformsIDs: []int32{td.Int31()},
		Screenshots:  []string{s.getImageURL()},
		Websites:     []string{s.getWebsiteURL()},
	}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: td.String(), UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().CreateGame(mock.Any(), mock.Any()).Return(gameID, nil)

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.CreateGame)))

	handler.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusCreated, s.httpResponse.Code)
	s.JSONEq(fmt.Sprintf(`{"id": %d}`, gameID), s.httpResponse.Body.String())
}

func (s *TestSuite) Test_CreateGame_DecodeError() {
	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewReader([]byte("{invalid json}")))

	s.provider.CreateGame(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
}

func (s *TestSuite) Test_CreateGame_MissingClaims() {
	requestData := api.CreateGameRequest{
		Name:         td.String(),
		Developer:    td.String(),
		ReleaseDate:  td.Date().Format("2006-01-02"),
		GenresIDs:    []int32{td.Int31()},
		LogoURL:      s.getImageURL(),
		Summary:      td.String(),
		PlatformsIDs: []int32{td.Int31()},
		Screenshots:  []string{s.getImageURL()},
		Websites:     []string{s.getWebsiteURL()},
	}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewReader(requestBody))

	s.provider.CreateGame(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_CreateGame_FacadeError() {
	role, authToken := td.String(), td.String()

	requestData := api.CreateGameRequest{
		Name:         td.String(),
		Developer:    td.String(),
		ReleaseDate:  td.Date().Format("2006-01-02"),
		GenresIDs:    []int32{td.Int31()},
		LogoURL:      s.getImageURL(),
		Summary:      td.String(),
		PlatformsIDs: []int32{td.Int31()},
		Screenshots:  []string{s.getImageURL()},
		Websites:     []string{s.getWebsiteURL()},
	}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: td.String(), UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().CreateGame(mock.Any(), mock.Any()).Return(int32(0), errors.New("db error"))

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.CreateGame)))

	handler.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) getImageURL() string {
	return fmt.Sprintf("https://ucarecdn.com/%s.jpg", td.String())
}

func (s *TestSuite) getWebsiteURL() string {
	var websites = []string{
		"https://gog.com/",
		"https://twitch.tv/",
		"https://youtube.com/",
	}
	return websites[td.Intn(len(websites))] + td.String()
}
