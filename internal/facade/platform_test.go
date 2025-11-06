package facade_test

import (
	"errors"
	"time"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	goredis "github.com/redis/go-redis/v9"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetPlatforms_Success() {
	platforms := []model.Platform{{
		ID:           td.Int32(),
		Name:         td.String(),
		Abbreviation: td.String(),
		IGDBID:       td.Int64(),
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetPlatforms(s.ctx).Return(platforms, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), platforms, time.Duration(0)).Return(nil)

	res, err := s.provider.GetPlatforms(s.ctx)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(platforms[0], res[0])
}

func (s *TestSuite) TestGetPlatforms_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetPlatforms(s.ctx).Return(nil, errors.New("new error"))

	res, err := s.provider.GetPlatforms(s.ctx)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetPlatformsMap_Success() {
	platforms := []model.Platform{{
		ID:           td.Int32(),
		Name:         td.String(),
		Abbreviation: td.String(),
		IGDBID:       td.Int64(),
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetPlatforms(s.ctx).Return(platforms, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), platforms, time.Duration(0)).Return(nil)

	res, err := s.provider.GetPlatformsMap(s.ctx)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(platforms[0], res[platforms[0].ID])
}

func (s *TestSuite) TestGetPlatformsMap_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetPlatforms(s.ctx).Return(nil, errors.New("new error"))

	res, err := s.provider.GetPlatformsMap(s.ctx)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetPlatformByID_Success() {
	platform := model.Platform{
		ID:           td.Int32(),
		Name:         td.String(),
		Abbreviation: td.String(),
		IGDBID:       td.Int64(),
	}

	s.storageMock.EXPECT().GetPlatformByID(s.ctx, platform.ID).Return(platform, nil)

	res, err := s.provider.GetPlatformByID(s.ctx, platform.ID)

	s.Require().NoError(err)
	s.Equal(platform.ID, res.ID)
	s.Equal(platform.Name, res.Name)
	s.Equal(platform.Abbreviation, res.Abbreviation)
	s.Equal(platform.IGDBID, res.IGDBID)
}

func (s *TestSuite) TestGetPlatformByID_NotFound() {
	platformID := td.Int32()

	s.storageMock.EXPECT().GetPlatformByID(s.ctx, platformID).Return(model.Platform{}, apperr.NewNotFoundError("platform", platformID))

	_, err := s.provider.GetPlatformByID(s.ctx, platformID)

	s.Require().Error(err)
	s.True(apperr.IsStatusCode(err, apperr.NotFound))
}

func (s *TestSuite) TestGetPlatformByID_Error() {
	platformID := td.Int32()

	s.storageMock.EXPECT().GetPlatformByID(s.ctx, platformID).Return(model.Platform{}, errors.New("new error"))

	_, err := s.provider.GetPlatformByID(s.ctx, platformID)

	s.Error(err)
}
