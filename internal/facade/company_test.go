package facade_test

import (
	"database/sql"
	"errors"
	"time"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	goredis "github.com/redis/go-redis/v9"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetCompanies_Success() {
	companies := []model.Company{{
		ID:   td.Int32(),
		Name: td.String(),
		IGDBID: sql.NullInt64{
			Int64: td.Int64(),
			Valid: td.Bool(),
		},
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetCompanies(s.ctx).Return(companies, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), companies, time.Duration(0)).Return(nil)

	res, err := s.provider.GetCompanies(s.ctx)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(companies[0], res[0])
}

func (s *TestSuite) TestGetCompanies_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetCompanies(s.ctx).Return(nil, errors.New("new error"))

	res, err := s.provider.GetCompanies(s.ctx)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetCompaniesMap_Success() {
	companies := []model.Company{{
		ID:   td.Int32(),
		Name: td.String(),
		IGDBID: sql.NullInt64{
			Int64: td.Int64(),
			Valid: td.Bool(),
		},
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetCompanies(s.ctx).Return(companies, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), companies, time.Duration(0)).Return(nil)

	res, err := s.provider.GetCompaniesMap(s.ctx)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(companies[0], res[companies[0].ID])
}

func (s *TestSuite) TestGetCompaniesMap_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetCompanies(s.ctx).Return(nil, errors.New("new error"))

	res, err := s.provider.GetCompaniesMap(s.ctx)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetTopCompanies_Success() {
	companies := []model.Company{{
		ID:   td.Int32(),
		Name: td.String(),
		IGDBID: sql.NullInt64{
			Int64: td.Int64(),
			Valid: td.Bool(),
		},
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetTopDevelopers(s.ctx, mock.Any()).Return(companies, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), companies, time.Duration(0)).Return(nil)

	res, err := s.provider.GetTopCompanies(s.ctx, model.CompanyTypeDeveloper, 5)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(companies[0], res[0])
}

func (s *TestSuite) TestGetTopCompanies_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetTopPublishers(s.ctx, mock.Any()).Return(nil, errors.New("new error"))

	res, err := s.provider.GetTopCompanies(s.ctx, model.CompanyTypePublisher, 5)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetCompanyByID_Success() {
	company := model.Company{
		ID:   td.Int32(),
		Name: td.String(),
	}

	s.storageMock.EXPECT().GetCompanyByID(s.ctx, company.ID).Return(company, nil)

	res, err := s.provider.GetCompanyByID(s.ctx, company.ID)

	s.Require().NoError(err)
	s.Equal(company, res)
}

func (s *TestSuite) TestGetCompanyByID_Error() {
	s.storageMock.EXPECT().GetCompanyByID(s.ctx, mock.Any()).Return(model.Company{}, errors.New("not found"))

	res, err := s.provider.GetCompanyByID(s.ctx, td.Int32())

	s.Require().Error(err)
	s.Equal(model.Company{}, res)
}

func (s *TestSuite) TestCompanyExistsInIGDB_Success_Found() {
	companyName := "Nintendo"

	s.igdbAPIClientMock.EXPECT().CompanyExists(s.ctx, companyName).Return(true, nil)

	exists, err := s.provider.CompanyExistsInIGDB(s.ctx, companyName)

	s.Require().NoError(err)
	s.True(exists)
}

func (s *TestSuite) TestCompanyExistsInIGDB_Success_NotFound() {
	companyName := "NonExistentCompany"

	s.igdbAPIClientMock.EXPECT().CompanyExists(s.ctx, companyName).Return(false, nil)

	exists, err := s.provider.CompanyExistsInIGDB(s.ctx, companyName)

	s.Require().NoError(err)
	s.False(exists)
}

func (s *TestSuite) TestCompanyExistsInIGDB_Error() {
	companyName := "SomeCompany"
	igdbErr := errors.New("IGDB API error")

	s.igdbAPIClientMock.EXPECT().CompanyExists(s.ctx, companyName).Return(false, igdbErr)

	exists, err := s.provider.CompanyExistsInIGDB(s.ctx, companyName)

	s.Require().Error(err)
	s.False(exists)
	s.Equal(igdbErr, err)
}

func (s *TestSuite) TestCompanyExistsInIGDB_EmptyCompanyName() {
	companyName := ""

	s.igdbAPIClientMock.EXPECT().CompanyExists(s.ctx, companyName).Return(false, nil)

	exists, err := s.provider.CompanyExistsInIGDB(s.ctx, companyName)

	s.Require().NoError(err)
	s.False(exists)
}

func (s *TestSuite) TestCompanyExistsInIGDB_MultipleCompanies() {
	companies := []string{"Sony", "Microsoft", "Nintendo", "Valve"}

	for _, company := range companies {
		s.igdbAPIClientMock.EXPECT().CompanyExists(s.ctx, company).Return(true, nil)
	}

	for _, company := range companies {
		exists, err := s.provider.CompanyExistsInIGDB(s.ctx, company)
		s.Require().NoError(err)
		s.True(exists)
	}
}
