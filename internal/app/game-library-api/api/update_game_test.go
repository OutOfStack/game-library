package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/go-chi/chi/v5"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_UpdateGame_Success() {
	gameID, authToken, publisher, role := td.Int31(), td.String(), td.String(), td.String()

	name, developer, releaseDate, summary, logoURL := td.String(), td.String(), td.Date().Format("2006-01-02"), td.String(), s.getImageURL()
	genresIDs, platformsIDs, screenshots, websites := []int32{td.Int31()}, []int32{td.Int31()}, []string{s.getImageURL()}, []string{s.getWebsiteURL()}
	requestData := api.UpdateGameRequest{
		Name:         &name,
		Developer:    &developer,
		ReleaseDate:  &releaseDate,
		GenresIDs:    &genresIDs,
		LogoURL:      &logoURL,
		Summary:      &summary,
		PlatformsIDs: &platformsIDs,
		Screenshots:  &screenshots,
		Websites:     &websites,
	}
	updateGame := model.UpdateGame{
		Name:         requestData.Name,
		Developer:    requestData.Developer,
		Publisher:    publisher,
		ReleaseDate:  requestData.ReleaseDate,
		GenresIDs:    requestData.GenresIDs,
		LogoURL:      requestData.LogoURL,
		Summary:      requestData.Summary,
		PlatformsIDs: requestData.PlatformsIDs,
		Screenshots:  requestData.Screenshots,
		Websites:     requestData.Websites,
	}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/games/%d", gameID), bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: publisher, UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().UpdateGame(mock.Any(), gameID, updateGame).Return(nil)

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.UpdateGame)))
	r := chi.NewRouter()
	r.Patch("/games/{id}", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusNoContent, s.httpResponse.Code)
}

func (s *TestSuite) Test_UpdateGame_InvalidID() {
	req := httptest.NewRequest(http.MethodPatch, "/games/-100", nil)

	r := chi.NewRouter()
	r.Patch("/games/{id}", s.provider.UpdateGame)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
	s.JSONEq(`{"error":"invalid id"}`, s.httpResponse.Body.String())
}

func (s *TestSuite) Test_UpdateGame_MissingClaims() {
	gameID := td.Int31()

	requestData := api.UpdateGameRequest{}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/games/%d", gameID), bytes.NewReader(requestBody))

	handler := http.HandlerFunc(s.provider.UpdateGame)
	r := chi.NewRouter()
	r.Patch("/games/{id}", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) Test_UpdateGame_FacadeError() {
	gameID, authToken, publisher, role := td.Int31(), td.String(), td.String(), td.String()

	requestData := api.UpdateGameRequest{}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/games/%d", gameID), bytes.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClient.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: publisher, UserRole: role}, nil)
	s.authClient.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().UpdateGame(mock.Any(), gameID, mock.Any()).Return(errors.New("new error"))

	authenticator := middleware.Authenticate(s.log, s.authClient)
	authorizer := middleware.Authorize(s.log, s.authClient, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.UpdateGame)))
	r := chi.NewRouter()
	r.Patch("/games/{id}", handler.ServeHTTP)

	r.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}
