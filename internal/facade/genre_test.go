package facade_test

import (
	"errors"
	"time"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	goredis "github.com/redis/go-redis/v9"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetGenres_Success() {
	genres := []model.Genre{{
		ID:     td.Int32(),
		Name:   td.String(),
		IGDBID: td.Int64(),
	}}

	s.storageMock.EXPECT().GetGenres(s.ctx).Return(genres, nil)
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), genres, time.Duration(0)).Return(nil)

	res, err := s.provider.GetGenres(s.ctx)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(genres[0], res[0])
}

func (s *TestSuite) TestGetGenres_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetGenres(s.ctx).Return(nil, errors.New("new error"))

	res, err := s.provider.GetGenres(s.ctx)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetGenresMap_Success() {
	genres := []model.Genre{{
		ID:     td.Int32(),
		Name:   td.String(),
		IGDBID: td.Int64(),
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetGenres(s.ctx).Return(genres, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), genres, time.Duration(0)).Return(nil)

	res, err := s.provider.GetGenresMap(s.ctx)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(genres[0], res[genres[0].ID])
}

func (s *TestSuite) TestGetGenresMap_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetGenres(s.ctx).Return(nil, errors.New("new error"))

	res, err := s.provider.GetGenresMap(s.ctx)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetTopGenres_Success() {
	genres := []model.Genre{{
		ID:     td.Int32(),
		Name:   td.String(),
		IGDBID: td.Int64(),
	}}

	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetTopGenres(s.ctx, mock.Any()).Return(genres, nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), genres, time.Duration(0)).Return(nil)

	res, err := s.provider.GetTopGenres(s.ctx, 5)

	s.Require().NoError(err)
	s.Len(res, 1)
	s.Equal(genres[0], res[0])
}

func (s *TestSuite) TestGetTopGenres_Error() {
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(errors.New("new error"))
	s.storageMock.EXPECT().GetTopGenres(s.ctx, mock.Any()).Return(nil, errors.New("new error"))

	res, err := s.provider.GetTopGenres(s.ctx, 5)

	s.Require().Error(err)
	s.Empty(res)
}

func (s *TestSuite) TestGetGenreByID_Success() {
	genre := model.Genre{
		ID:     td.Int32(),
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	s.storageMock.EXPECT().GetGenreByID(s.ctx, genre.ID).Return(genre, nil)

	res, err := s.provider.GetGenreByID(s.ctx, genre.ID)

	s.Require().NoError(err)
	s.Equal(genre, res)
}

func (s *TestSuite) TestGetGenreByID_Error() {
	s.storageMock.EXPECT().GetGenreByID(s.ctx, mock.Any()).Return(model.Genre{}, errors.New("not found"))

	res, err := s.provider.GetGenreByID(s.ctx, td.Int32())

	s.Require().Error(err)
	s.Equal(model.Genre{}, res)
}
