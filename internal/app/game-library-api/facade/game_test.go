package facade_test

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/facade"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/OutOfStack/game-library/pkg/types"
	goredis "github.com/redis/go-redis/v9"
	mock "go.uber.org/mock/gomock"
)

func (s *TestSuite) TestGetGames_Success() {
	games := []model.Game{{
		ID:            td.Int32(),
		Name:          td.String(),
		DevelopersIDs: []int32{td.Int32(), td.Int32()},
		PublishersIDs: []int32{td.Int32(), td.Int32()},
		ReleaseDate:   types.DateOf(td.Date()),
		GenresIDs:     []int32{td.Int32(), td.Int32()},
		LogoURL:       td.String(),
		Rating:        td.Float64n(5),
		Summary:       td.String(),
		Slug:          td.String(),
		PlatformsIDs:  []int32{td.Int32(), td.Int32()},
		Screenshots:   []string{td.String(), td.String()},
		Websites:      []string{td.String(), td.String()},
		IGDBRating:    td.Float64n(100),
		IGDBID:        td.Int64(),
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
	s.storageMock.EXPECT().GetGamesCount(s.ctx, mock.Any()).Return(uint64(0), errors.New("count error"))
	s.redisClientMock.EXPECT().GetStruct(s.ctx, mock.Any(), mock.Any()).Return(goredis.Nil).Times(2)

	res, cnt, err := s.provider.GetGames(s.ctx, 1, 10, model.GamesFilter{})

	s.Error(err)
	s.Len(res, 0)
	s.Equal(uint64(0), cnt)
}

func (s *TestSuite) TestGetGameByID_Success() {
	game := model.Game{
		ID:            td.Int32(),
		Name:          td.String(),
		DevelopersIDs: []int32{td.Int32(), td.Int32()},
		PublishersIDs: []int32{td.Int32(), td.Int32()},
		ReleaseDate:   types.DateOf(td.Date()),
		GenresIDs:     []int32{td.Int32(), td.Int32()},
		LogoURL:       td.String(),
		Rating:        td.Float64n(5),
		Summary:       td.String(),
		Slug:          td.String(),
		PlatformsIDs:  []int32{td.Int32(), td.Int32()},
		Screenshots:   []string{td.String(), td.String()},
		Websites:      []string{td.String(), td.String()},
		IGDBRating:    td.Float64n(100),
		IGDBID:        td.Int64(),
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
		Name:         td.String(),
		Developer:    td.String(),
		Publisher:    td.String(),
		ReleaseDate:  td.Date().String(),
		GenresIDs:    []int32{td.Int32(), td.Int32()},
		LogoURL:      td.String(),
		Summary:      td.String(),
		Slug:         td.String(),
		PlatformsIDs: []int32{td.Int32(), td.Int32()},
		Screenshots:  []string{td.String(), td.String()},
		Websites:     []string{td.String(), td.String()},
	}

	createGameData := model.CreateGameData{
		Name:             createGame.Name,
		DevelopersIDs:    []int32{developerID},
		PublishersIDs:    []int32{publisherID},
		ReleaseDate:      createGame.ReleaseDate,
		GenresIDs:        createGame.GenresIDs,
		LogoURL:          createGame.LogoURL,
		Summary:          createGame.Summary,
		Slug:             createGame.Slug,
		PlatformsIDs:     createGame.PlatformsIDs,
		Screenshots:      createGame.Screenshots,
		Websites:         createGame.Websites,
		ModerationStatus: model.ModerationStatusCheck,
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Millisecond)

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Developer).Return(int32(0), nil)
	s.storageMock.EXPECT().CreateCompany(s.ctx, model.Company{Name: createGame.Developer}).Return(developerID, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().GetPublisherGamesCount(s.ctx, publisherID, startOfMonth, endOfMonth).Return(1, nil)
	s.storageMock.EXPECT().CreateGame(s.ctx, createGameData).Return(gameID, nil)
	s.storageMock.EXPECT().GetGameTrendingData(mock.Any(), gameID).Return(model.GameTrendingData{}, nil).AnyTimes()
	s.storageMock.EXPECT().UpdateGameTrendingIndex(mock.Any(), gameID, mock.Any()).Return(nil).AnyTimes()

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
		Developer: td.String(),
		Publisher: td.String(),
	}

	createGameData := model.CreateGameData{
		DevelopersIDs:    []int32{developerID},
		PublishersIDs:    []int32{publisherID},
		ModerationStatus: model.ModerationStatusCheck,
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Developer).Return(developerID, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().GetPublisherGamesCount(s.ctx, publisherID, mock.Any(), mock.Any()).Return(1, nil)
	s.storageMock.EXPECT().CreateGame(s.ctx, createGameData).Return(int32(0), errors.New("new error"))

	id, err := s.provider.CreateGame(s.ctx, createGame)

	s.Error(err)
	s.Equal(int32(0), id)
}

func (s *TestSuite) TestCreateGame_MonthlyLimitReached() {
	developerID, publisherID := td.Int32(), td.Int32()
	createGame := model.CreateGame{
		Developer: td.String(),
		Publisher: td.String(),
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Developer).Return(developerID, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, createGame.Publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().GetPublisherGamesCount(s.ctx, publisherID, mock.Any(), mock.Any()).Return(facade.MaxGamesPerPublisherPerMonth, nil)

	id, err := s.provider.CreateGame(s.ctx, createGame)

	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("publishing monthly limit of %d reached", facade.MaxGamesPerPublisherPerMonth))
	s.Equal(int32(0), id)
}

func (s *TestSuite) TestUpdateGame_Success() {
	game := model.Game{
		ID:            td.Int32(),
		PublishersIDs: []int32{td.Int32()},
	}
	updateGame := model.UpdateGame{
		Publisher: td.String(),
	}
	updateGameData := model.UpdateGameData{
		PublishersIDs:    game.PublishersIDs,
		ModerationStatus: model.ModerationStatusRecheck,
	}

	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, updateGame.Publisher).Return(game.PublishersIDs[0], nil)
	s.storageMock.EXPECT().UpdateGame(s.ctx, game.ID, updateGameData).Return(nil)
	s.storageMock.EXPECT().GetGameTrendingData(mock.Any(), game.ID).Return(model.GameTrendingData{}, nil).AnyTimes()
	s.storageMock.EXPECT().UpdateGameTrendingIndex(mock.Any(), game.ID, mock.Any()).Return(nil).AnyTimes()

	s.redisClientMock.EXPECT().DeleteByMatch(mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().GetStruct(mock.Any(), mock.Any(), mock.Any()).Return(nil).AnyTimes()
	s.redisClientMock.EXPECT().Delete(mock.Any(), mock.Any()).Return(nil).AnyTimes()

	err := s.provider.UpdateGame(s.ctx, game.ID, updateGame)

	s.NoError(err)
}

func (s *TestSuite) TestUpdateGame_Forbidden() {
	game := model.Game{
		ID: td.Int32(),
	}
	updateGame := model.UpdateGame{
		Publisher: td.String(),
	}

	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, updateGame.Publisher).Return(td.Int32(), nil)

	err := s.provider.UpdateGame(s.ctx, game.ID, updateGame)

	s.Error(err)
	s.True(apperr.IsStatusCode(err, http.StatusForbidden))
}

func (s *TestSuite) TestUpdateGame_Error() {
	gameID := td.Int32()
	updateGame := model.UpdateGame{
		Publisher: td.String(),
	}

	s.storageMock.EXPECT().GetGameByID(s.ctx, gameID).Return(model.Game{}, errors.New("new error"))

	err := s.provider.UpdateGame(s.ctx, gameID, updateGame)

	s.Error(err)
}

func (s *TestSuite) TestDeleteGame_Success() {
	publisher := td.String()
	game := model.Game{
		ID:            td.Int32(),
		PublishersIDs: []int32{td.Int32()},
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(game.PublishersIDs[0], nil)
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
		ID:            td.Int32(),
		PublishersIDs: []int32{td.Int32()},
	}

	s.storageMock.EXPECT().GetCompanyIDByName(s.ctx, publisher).Return(game.PublishersIDs[0], nil)
	s.storageMock.EXPECT().GetGameByID(s.ctx, game.ID).Return(game, nil)
	s.storageMock.EXPECT().DeleteGame(s.ctx, game.ID).Return(errors.New("new error"))

	err := s.provider.DeleteGame(s.ctx, game.ID, publisher)

	s.Error(err)
}

func (s *TestSuite) TestUpdateGameTrendingIndex_Success() {
	gameID := td.Int32()
	trendingData := model.GameTrendingData{
		Year:            2000 + td.Intn(25),
		Month:           1 + td.Intn(12),
		IGDBRating:      td.Float64n(100),
		IGDBRatingCount: 100 + int32(td.Intn(100)),
		Rating:          td.Float64n(5),
		RatingCount:     td.Int31(),
	}

	s.storageMock.EXPECT().GetGameTrendingData(s.ctx, gameID).Return(trendingData, nil)
	s.storageMock.EXPECT().UpdateGameTrendingIndex(s.ctx, gameID, mock.Any()).Return(nil)

	err := s.provider.UpdateGameTrendingIndex(s.ctx, gameID)

	s.NoError(err)
}

func (s *TestSuite) TestUpdateGameTrendingIndex_GetDataError() {
	gameID := td.Int32()

	s.storageMock.EXPECT().GetGameTrendingData(s.ctx, gameID).Return(model.GameTrendingData{}, errors.New("get data error"))

	err := s.provider.UpdateGameTrendingIndex(s.ctx, gameID)

	s.Error(err)
	s.Contains(err.Error(), "get data error")
}

func (s *TestSuite) TestUpdateGameTrendingIndex_UpdateError() {
	gameID := td.Int32()
	trendingData := model.GameTrendingData{
		Year:            2000 + td.Intn(25),
		Month:           1 + td.Intn(12),
		IGDBRating:      td.Float64n(100),
		IGDBRatingCount: 100 + int32(td.Intn(100)),
		Rating:          td.Float64n(5),
		RatingCount:     td.Int31(),
	}

	s.storageMock.EXPECT().GetGameTrendingData(s.ctx, gameID).Return(trendingData, nil)
	s.storageMock.EXPECT().UpdateGameTrendingIndex(s.ctx, gameID, mock.Any()).Return(errors.New("update error"))

	err := s.provider.UpdateGameTrendingIndex(s.ctx, gameID)

	s.Error(err)
	s.Contains(err.Error(), "update error")
}
