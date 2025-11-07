package facade_test

import (
	"errors"
	"time"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	goredis "github.com/redis/go-redis/v9"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestRateGame_Success() {
	gameID, userID, rating := td.Int32(), td.String(), uint8(td.Intn(5)+1)

	s.storageMock.EXPECT().GetGameByID(s.ctx, gameID).Return(model.Game{}, nil).Times(1)
	s.storageMock.EXPECT().AddRating(s.ctx, model.CreateRating{
		Rating: rating,
		UserID: userID,
		GameID: gameID,
	}).Return(nil).Times(1)

	s.storageMock.EXPECT().UpdateGameRating(mock.Any(), gameID).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().Delete(mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().GetStruct(mock.Any(), mock.Any(), mock.Any()).Return(nil).AnyTimes()

	err := s.provider.RateGame(s.ctx, gameID, userID, rating)

	s.Require().NoError(err)
}

func (s *TestSuite) TestRateGame_Delete_Success() {
	gameID, userID, rating := td.Int32(), td.String(), uint8(0)

	s.storageMock.EXPECT().GetGameByID(s.ctx, gameID).Return(model.Game{}, nil).Times(1)
	s.storageMock.EXPECT().RemoveRating(s.ctx, model.RemoveRating{
		UserID: userID,
		GameID: gameID,
	}).Return(nil).Times(1)

	s.storageMock.EXPECT().UpdateGameRating(mock.Any(), gameID).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().Delete(mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().GetStruct(mock.Any(), mock.Any(), mock.Any()).Return(nil).AnyTimes()

	err := s.provider.RateGame(s.ctx, gameID, userID, rating)

	s.Require().NoError(err)
}

func (s *TestSuite) TestRateGame_Error() {
	gameID, userID, rating := td.Int32(), td.String(), uint8(0)

	s.storageMock.EXPECT().GetGameByID(s.ctx, gameID).Return(model.Game{}, errors.New("new error")).Times(1)

	err := s.provider.RateGame(s.ctx, gameID, userID, rating)

	s.Error(err)
}

func (s *TestSuite) TestGetUserRatings_Success() {
	userID := td.String()
	m := map[int32]uint8{td.Int32(): td.Uint8()}

	s.storageMock.EXPECT().GetUserRatings(s.ctx, userID).Return(m, nil).Times(1)
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil).Times(1)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), m, time.Duration(0)).Return(nil).Times(1)

	res, err := s.provider.GetUserRatings(s.ctx, userID)

	s.Require().NoError(err)
	s.Equal(m, res)
}

func (s *TestSuite) TestGetUserRatings_Error() {
	userID := td.String()

	s.storageMock.EXPECT().GetUserRatings(s.ctx, userID).Return(nil, errors.New("new error")).Times(1)
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil).Times(1)

	res, err := s.provider.GetUserRatings(s.ctx, userID)

	s.Require().Error(err)
	s.Nil(res)
}
