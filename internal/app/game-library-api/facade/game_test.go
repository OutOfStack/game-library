package facade_test

import (
	"errors"
	"net/http"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/OutOfStack/game-library/pkg/types"
	goredis "github.com/redis/go-redis/v9"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetGames_Success() {
	games := []model.Game{{
		ID:          td.Int32(),
		Name:        td.String(),
		Developers:  []int32{td.Int32(), td.Int32()},
		Publishers:  []int32{td.Int32(), td.Int32()},
		ReleaseDate: types.DateOf(td.Date()),
		Genres:      []int32{td.Int32(), td.Int32()},
		LogoURL:     td.String(),
		Rating:      td.Float64(),
		Summary:     td.String(),
		Slug:        td.String(),
		Platforms:   []int32{td.Int32(), td.Int32()},
		Screenshots: []string{td.String(), td.String()},
		Websites:    []string{td.String(), td.String()},
		IGDBRating:  td.Float64(),
		IGDBID:      td.Int64(),
	}}
	var count = td.Uint64()

	s.storageMock.EXPECT().GetGames(s.ctx, mock.Any(), mock.Any(), mock.Any()).Return(games, nil)
	s.storageMock.EXPECT().GetGamesCount(s.ctx, mock.Any()).Return(count, nil)
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil).Times(2)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), games, time.Duration(0)).Return(nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), count, time.Duration(0)).Return(nil)

	res, cnt, err := s.provider.GetGames(s.ctx, 1, 10, model.GamesFilter{})

	s.NoError(err)
	s.Len(res, 1)
	s.Equal(games[0], res[0])
	s.Equal(count, cnt)
}

func (s *TestSuite) TestGetGames_Error() {
	s.storageMock.EXPECT().GetGames(s.ctx, mock.Any(), mock.Any(), mock.Any()).Return(nil, errors.New("new error"))
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)

	res, cnt, err := s.provider.GetGames(s.ctx, 1, 10, model.GamesFilter{})

	s.Error(err)
	s.Len(res, 0)
	s.Equal(uint64(0), cnt)
}

func (s *TestSuite) TestGetGameByID_Success() {
	game := model.Game{
		ID:          td.Int32(),
		Name:        td.String(),
		Developers:  []int32{td.Int32(), td.Int32()},
		Publishers:  []int32{td.Int32(), td.Int32()},
		ReleaseDate: types.DateOf(td.Date()),
		Genres:      []int32{td.Int32(), td.Int32()},
		LogoURL:     td.String(),
		Rating:      td.Float64(),
		Summary:     td.String(),
		Slug:        td.String(),
		Platforms:   []int32{td.Int32(), td.Int32()},
		Screenshots: []string{td.String(), td.String()},
		Websites:    []string{td.String(), td.String()},
		IGDBRating:  td.Float64(),
		IGDBID:      td.Int64(),
	}

	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(s.ctx, mock.Any(), game, time.Duration(0)).Return(nil)

	res, err := s.provider.GetGameByID(s.ctx, game.ID)

	s.NoError(err)
	s.Equal(game, res)
}

func (s *TestSuite) TestGetGameByID_Error() {
	s.storageMock.EXPECT().GetGameByID(s.ctx, mock.Any()).Return(model.Game{}, errors.New("not found"))
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil)

	res, err := s.provider.GetGameByID(s.ctx, td.Int32())

	s.NotNil(err)
	s.Equal(model.Game{}, res)
}

func (s *TestSuite) TestCreateGame_Success() {
	gameID := td.Int31()
	developerID, publisherID := td.Int31(), td.Int31()
	createGame := model.CreateGame{
		Name:          td.String(),
		Developer:     td.String(),
		DevelopersIDs: []int32{developerID},
		Publisher:     td.String(),
		PublishersIDs: []int32{publisherID},
		ReleaseDate:   td.Date().String(),
		Genres:        []int32{td.Int32(), td.Int32()},
		LogoURL:       td.String(),
		Summary:       td.String(),
		Slug:          td.String(),
		Platforms:     []int32{td.Int32(), td.Int32()},
		Screenshots:   []string{td.String(), td.String()},
		Websites:      []string{td.String(), td.String()},
		IGDBRating:    td.Float64(),
		IGDBID:        td.Int64(),
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Developer).Return(int32(0), nil)
	s.storageMock.EXPECT().CreateCompany(s.ctx, model.Company{Name: createGame.Developer}).Return(developerID, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().CreateGame(s.ctx, createGame).Return(gameID, nil)

	s.redisClientMock.EXPECT().DeleteByMatch(mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().GetStruct(mock.Any(), mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().Delete(mock.Any(), mock.Any()).Return(nil).AnyTimes()

	id, err := s.provider.CreateGame(s.ctx, createGame)

	s.NoError(err)
	s.Equal(gameID, id)
}

func (s *TestSuite) TestCreateGame_Error() {
	developerID, publisherID := td.Int32(), td.Int32()
	createGame := model.CreateGame{
		Developer:     td.String(),
		DevelopersIDs: []int32{developerID},
		Publisher:     td.String(),
		PublishersIDs: []int32{publisherID},
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Developer).Return(developerID, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().CreateGame(s.ctx, createGame).Return(int32(0), errors.New("new error"))

	id, err := s.provider.CreateGame(s.ctx, createGame)

	s.Error(err)
	s.Equal(int32(0), id)
}

func (s *TestSuite) TestUpdateGame_Success() {
	game := model.Game{
		ID:         td.Int32(),
		Publishers: []int32{td.Int32()},
	}
	publisher := td.String()
	updatedGame := model.UpdatedGame{}

	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(game.Publishers[0], nil)
	s.storageMock.EXPECT().UpdateGame(s.ctx, game.ID, mock.Any()).Return(nil)

	s.redisClientMock.EXPECT().DeleteByMatch(mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().GetStruct(mock.Any(), mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().Delete(mock.Any(), mock.Any()).Return(nil).AnyTimes()

	err := s.provider.UpdateGame(s.ctx, game.ID, publisher, updatedGame)

	s.NoError(err)
}

func (s *TestSuite) TestUpdateGame_Forbidden() {
	game := model.Game{
		ID: td.Int32(),
	}
	publisher := td.String()
	updatedGame := model.UpdatedGame{}

	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(td.Int32(), nil)

	err := s.provider.UpdateGame(s.ctx, game.ID, publisher, updatedGame)

	s.Error(err)
	s.True(apperr.IsStatusCode(err, http.StatusForbidden))
}

func (s *TestSuite) TestUpdateGame_Error() {
	gameID := td.Int32()
	publisher := td.String()
	updatedGame := model.UpdatedGame{}

	s.storageMock.EXPECT().GetGameByID(s.ctx, gameID).Return(model.Game{}, errors.New("new error"))

	err := s.provider.UpdateGame(s.ctx, gameID, publisher, updatedGame)

	s.Error(err)
}

func (s *TestSuite) TestDeleteGame_Success() {
	publisher := td.String()
	game := model.Game{
		ID:         td.Int32(),
		Publishers: []int32{td.Int32()},
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(game.Publishers[0], nil)
	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().DeleteGame(s.ctx, game.ID).Return(nil)

	s.redisClientMock.EXPECT().DeleteByMatch(mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().Delete(mock.Any(), mock.Any()).Return(nil).AnyTimes()

	err := s.provider.DeleteGame(s.ctx, game.ID, publisher)

	s.NoError(err)
}

func (s *TestSuite) TestDeleteGame_Forbidden() {
	publisher := td.String()
	game := model.Game{
		ID: td.Int32(),
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(td.Int32(), nil)
	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)

	err := s.provider.DeleteGame(s.ctx, game.ID, publisher)

	s.Error(err)
	s.True(apperr.IsStatusCode(err, http.StatusForbidden))
}

func (s *TestSuite) TestDeleteGame_Error() {
	publisher := td.String()
	game := model.Game{
		ID:         td.Int32(),
		Publishers: []int32{td.Int32()},
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(game.Publishers[0], nil)
	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().DeleteGame(s.ctx, game.ID).Return(errors.New("new error"))

	err := s.provider.DeleteGame(s.ctx, game.ID, publisher)

	s.Error(err)
}
