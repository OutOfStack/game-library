package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) Test_UploadGameImages_Success() {
	userName, role, authToken := td.String(), td.String(), td.String()

	// create a multipart form buffer
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	coverWriter, _ := w.CreateFormFile("cover", "cover.jpg")
	coverContent := td.Bytes()
	_, err := io.Copy(coverWriter, bytes.NewReader(coverContent))
	s.NoError(err)

	screenshotWriter, _ := w.CreateFormFile("screenshots", "screenshot1.jpg")
	screenshotContent := td.Bytes()
	_, err = io.Copy(screenshotWriter, bytes.NewReader(screenshotContent))
	s.NoError(err)

	if cErr := w.Close(); cErr != nil {
		s.T().Log(cErr)
	}

	// create request with multipart form data
	req := httptest.NewRequest(http.MethodPost, "/games/images", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+authToken)

	expectedFiles := []model.File{
		{
			FileName: "cover.jpg",
			FileID:   td.String(),
			FileURL:  td.String(),
			Type:     "cover",
		},
		{
			FileName: "screenshot1.jpg",
			FileID:   td.String(),
			FileURL:  td.String(),
			Type:     "screenshot",
		},
	}

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: userName, UserRole: role}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().UploadGameImages(mock.Any(), mock.Any(), mock.Any(), userName).Return(expectedFiles, nil)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.UploadGameImages)))

	handler.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusCreated, s.httpResponse.Code)

	var response api.UploadImagesResponse
	err = json.NewDecoder(s.httpResponse.Body).Decode(&response)
	s.NoError(err)
	s.Len(response.Files, len(expectedFiles))
	for i, file := range response.Files {
		s.Equal(expectedFiles[i].FileName, file.FileName)
		s.Equal(expectedFiles[i].FileID, file.FileID)
		s.Equal(expectedFiles[i].FileURL, file.FileURL)
		s.Equal(expectedFiles[i].Type, file.Type)
	}
}

func (s *TestSuite) Test_UploadGameImages_ParseFormError() {
	role, authToken := td.String(), td.String()

	req := httptest.NewRequest(http.MethodPost, "/games/images", bytes.NewReader([]byte("invalid form data")))
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: td.String(), UserRole: role}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.UploadGameImages)))

	handler.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
	s.JSONEq(`{"error": "failed to parse form"}`, s.httpResponse.Body.String())
}

func (s *TestSuite) Test_UploadGameImages_MissingClaims() {
	req := httptest.NewRequest(http.MethodPost, "/games/images", nil)

	s.provider.UploadGameImages(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
	s.JSONEq(`{"error": "Internal Server Error"}`, s.httpResponse.Body.String())
}

func (s *TestSuite) Test_UploadGameImages_PublishingMonthlyLimitReached() {
	userName, role, authToken := td.String(), td.String(), td.String()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// add a test file
	fileWriter, _ := w.CreateFormFile("cover", "test.jpg")
	_, err := io.Copy(fileWriter, bytes.NewReader(td.Bytes()))
	s.NoError(err)

	if cErr := w.Close(); cErr != nil {
		s.T().Log(cErr)
	}

	req := httptest.NewRequest(http.MethodPost, "/games/images", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: userName, UserRole: role}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().UploadGameImages(mock.Any(), mock.Any(), mock.Any(), userName).Return(nil, apperr.NewTooManyRequestsError("game", "publishing monthly limit reached"))

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.UploadGameImages)))

	handler.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusTooManyRequests, s.httpResponse.Code)
	s.JSONEq(`{"error": "too many requests on game: publishing monthly limit reached"}`, s.httpResponse.Body.String())
}

func (s *TestSuite) Test_UploadGameImages_FacadeError() {
	userName, role, authToken := td.String(), td.String(), td.String()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// add a test file
	fileWriter, _ := w.CreateFormFile("cover", "test.jpg")
	_, err := io.Copy(fileWriter, bytes.NewReader(td.Bytes()))
	s.NoError(err)

	if cErr := w.Close(); cErr != nil {
		s.T().Log(cErr)
	}

	req := httptest.NewRequest(http.MethodPost, "/games/images", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+authToken)

	s.authClientMock.EXPECT().ParseToken(mock.Any()).Return(&auth.Claims{Name: userName, UserRole: role}, nil)
	s.authClientMock.EXPECT().Verify(mock.Any(), authToken).Return(nil)
	s.gameFacadeMock.EXPECT().UploadGameImages(mock.Any(), mock.Any(), mock.Any(), userName).Return(nil, errors.New("upload error"))

	authenticator := middleware.Authenticate(s.log, s.authClientMock)
	authorizer := middleware.Authorize(s.log, s.authClientMock, role)
	handler := authenticator(authorizer(http.HandlerFunc(s.provider.UploadGameImages)))

	handler.ServeHTTP(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
	s.JSONEq(`{"error": "Internal Server Error"}`, s.httpResponse.Body.String())
}
