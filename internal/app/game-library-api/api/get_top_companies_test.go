package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestTopGetCompanies_Success() {
	companies := []string{
		model.CompanyTypePublisher,
		model.CompanyTypeDeveloper,
	}
	companyType := companies[td.Intn(len(companies))]
	req := httptest.NewRequest(http.MethodGet, "/company/top/?type="+companyType, nil)

	s.gameFacadeMock.EXPECT().GetTopCompanies(mock.Any(), companyType, int64(10)).Return([]model.Company{
		{ID: 1, Name: "Developer 1"},
		{ID: 2, Name: "Publisher 2"},
	}, nil)

	s.provider.GetTopCompanies(s.httpResponse, req)

	s.Equal(http.StatusOK, s.httpResponse.Code)
	s.JSONEq(`[
			{"id": 1, "name": "Developer 1"},
			{"id": 2, "name": "Publisher 2"}
		]`, s.httpResponse.Body.String())
}

func (s *TestSuite) TestTopGetCompanies_Error() {
	companies := []string{
		model.CompanyTypePublisher,
		model.CompanyTypeDeveloper,
	}
	companyType := companies[td.Intn(len(companies))]
	req := httptest.NewRequest(http.MethodGet, "/company/top/?type="+companyType, nil)

	s.gameFacadeMock.EXPECT().GetTopCompanies(mock.Any(), companyType, int64(10)).Return(nil, errors.New("new error"))

	s.provider.GetTopCompanies(s.httpResponse, req)

	s.Equal(http.StatusInternalServerError, s.httpResponse.Code)
}

func (s *TestSuite) TestTopGetCompanies_NoTypeProvided() {
	s.provider.GetTopCompanies(s.httpResponse, s.httpRequest)

	s.Equal(http.StatusBadRequest, s.httpResponse.Code)
}
