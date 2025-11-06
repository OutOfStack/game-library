package infoapi_test

import (
	"errors"

	pb "github.com/OutOfStack/game-library/pkg/proto/infoapi"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TestSuite) Test_CompanyExists_Success_Found() {
	companyName := "Sony"
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()

	s.gameFacadeMock.EXPECT().
		CompanyExistsInIGDB(gomock.Any(), companyName).
		Return(true, nil)

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().True(resp.GetExists())
}

func (s *TestSuite) Test_CompanyExists_Success_NotFound() {
	companyName := "NonExistentCompany"
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()

	s.gameFacadeMock.EXPECT().
		CompanyExistsInIGDB(gomock.Any(), companyName).
		Return(false, nil)

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().False(resp.GetExists())
}

func (s *TestSuite) Test_CompanyExists_EmptyCompanyName() {
	var companyName string
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().Error(err)
	s.Require().Nil(resp)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("empty company name", st.Message())
}

func (s *TestSuite) Test_CompanyExists_WhitespaceOnlyCompanyName() {
	companyName := "   \t\n  "
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
	s.Equal("empty company name", st.Message())
}

func (s *TestSuite) Test_CompanyExists_FacadeError() {
	companyName := "SomeCompany"
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()
	facadeErr := errors.New("database error")

	s.gameFacadeMock.EXPECT().
		CompanyExistsInIGDB(gomock.Any(), companyName).
		Return(false, facadeErr)

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().Error(err)
	s.Nil(resp)
	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.Internal, st.Code())
	s.Equal("failed to check company existence", st.Message())
}

func (s *TestSuite) Test_CompanyExists_CaseInsensitive() {
	companyName := "NINTENDO"
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()

	s.gameFacadeMock.EXPECT().
		CompanyExistsInIGDB(gomock.Any(), companyName).
		Return(true, nil)

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().True(resp.GetExists())
}

func (s *TestSuite) Test_CompanyExists_CompanyNameWithSpaces() {
	companyName := "  Ubisoft  "
	req := pb.CompanyExistsRequest_builder{
		CompanyName: &companyName,
	}.Build()

	s.gameFacadeMock.EXPECT().
		CompanyExistsInIGDB(gomock.Any(), companyName).
		Return(true, nil)

	resp, err := s.service.CompanyExists(s.T().Context(), req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().True(resp.GetExists())
}
