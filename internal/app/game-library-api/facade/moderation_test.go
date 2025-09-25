package facade_test

import (
	"errors"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestCreateModerationRecord_Success() {
	gameID := td.Int31()
	game := model.Game{
		ID:            gameID,
		Name:          td.String(),
		DevelopersIDs: []int32{td.Int31()},
		PublishersIDs: []int32{td.Int31()},
		GenresIDs:     []int32{td.Int31()},
		LogoURL:       td.URL(),
		Summary:       td.String(),
		Slug:          td.String(),
		Screenshots:   []string{td.URL()},
		Websites:      []string{td.URL()},
	}
	companies := map[int32]model.Company{
		game.DevelopersIDs[0]: {Name: td.String()},
		game.PublishersIDs[0]: {Name: td.String()},
	}
	genres := map[int32]model.Genre{
		game.GenresIDs[0]: {Name: td.String()},
	}
	moderationID := td.Int31()

	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "companies", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "companies", gomock.Any(), gomock.Any()).Return(nil)
	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "genres", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "genres", gomock.Any(), gomock.Any()).Return(nil)
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return([]model.Company{
		{ID: game.DevelopersIDs[0], Name: companies[game.DevelopersIDs[0]].Name},
		{ID: game.PublishersIDs[0], Name: companies[game.PublishersIDs[0]].Name},
	}, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return([]model.Genre{
		{ID: game.GenresIDs[0], Name: genres[game.GenresIDs[0]].Name},
	}, nil)
	s.storageMock.EXPECT().CreateModerationRecord(gomock.Any(), gomock.Any()).Return(moderationID, nil)
	s.storageMock.EXPECT().UpdateGameModerationID(gomock.Any(), gameID, moderationID).Return(nil)

	id, err := s.provider.CreateModerationRecord(s.T().Context(), gameID)

	s.Require().NoError(err)
	s.Require().Equal(moderationID, id)
}

func (s *TestSuite) TestCreateModerationRecord_GetGameError() {
	gameID := td.Int31()
	getGameErr := errors.New("get game error")

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(model.Game{}, getGameErr)

	_, err := s.provider.CreateModerationRecord(s.T().Context(), gameID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "get game")
}

func (s *TestSuite) TestCreateModerationRecord_GetCompaniesError() {
	gameID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{td.Int31()},
	}
	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "companies", gomock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return(nil, errors.New(""))

	_, err := s.provider.CreateModerationRecord(s.T().Context(), gameID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "get game")
	s.Require().Contains(err.Error(), "moderation data")
}

func (s *TestSuite) TestCreateModerationRecord_GetGenresError() {
	gameID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{td.Int31()},
	}
	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "companies", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "companies", gomock.Any(), gomock.Any()).Return(nil)
	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "genres", gomock.Any()).Return(goredis.Nil)
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return([]model.Company{{ID: game.PublishersIDs[0], Name: td.String()}}, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return(nil, errors.New(""))

	_, err := s.provider.CreateModerationRecord(s.T().Context(), gameID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "get game")
	s.Require().Contains(err.Error(), "moderation data")
}

func (s *TestSuite) TestCreateModerationRecord_CreateModerationError() {
	gameID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{td.Int31()},
	}
	createModerationErr := errors.New("create moderation error")

	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "companies", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "companies", gomock.Any(), gomock.Any()).Return(nil)
	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "genres", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "genres", gomock.Any(), gomock.Any()).Return(nil)
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return([]model.Company{{ID: game.PublishersIDs[0], Name: td.String()}}, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return([]model.Genre{}, nil)
	s.storageMock.EXPECT().CreateModerationRecord(gomock.Any(), gomock.Any()).Return(int32(0), createModerationErr)

	_, err := s.provider.CreateModerationRecord(s.T().Context(), gameID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "create moderation record")
}

func (s *TestSuite) TestCreateModerationRecord_UpdateModerationIDError() {
	gameID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{td.Int31()},
	}
	moderationID := td.Int31()
	updateModerationIDErr := errors.New("update moderation id error")

	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "companies", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "companies", gomock.Any(), gomock.Any()).Return(nil)
	s.redisClientMock.EXPECT().GetStruct(gomock.Any(), "genres", gomock.Any()).Return(goredis.Nil)
	s.redisClientMock.EXPECT().SetStruct(gomock.Any(), "genres", gomock.Any(), gomock.Any()).Return(nil)
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return([]model.Company{{ID: game.PublishersIDs[0], Name: td.String()}}, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return([]model.Genre{}, nil)
	s.storageMock.EXPECT().CreateModerationRecord(gomock.Any(), gomock.Any()).Return(moderationID, nil)
	s.storageMock.EXPECT().UpdateGameModerationID(gomock.Any(), gameID, moderationID).Return(updateModerationIDErr)

	_, err := s.provider.CreateModerationRecord(s.T().Context(), gameID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "update game")
	s.Require().Contains(err.Error(), "moderation id")
}

func (s *TestSuite) TestGetGameModerations_Success() {
	gameID := td.Int31()
	publisher := td.String()
	publisherID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{publisherID},
	}
	moderations := []model.Moderation{
		{ID: td.Int31(), GameID: gameID},
		{ID: td.Int31(), GameID: gameID},
	}

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(gomock.Any(), publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().GetModerationRecordsByGameID(gomock.Any(), gameID).Return(moderations, nil)

	result, err := s.provider.GetGameModerations(s.T().Context(), gameID, publisher)

	s.Require().NoError(err)
	s.Require().Equal(moderations, result)
}

func (s *TestSuite) TestGetGameModerations_GetGameError() {
	gameID := td.Int31()
	publisher := td.String()
	getGameErr := errors.New("get game error")

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(model.Game{}, getGameErr)

	_, err := s.provider.GetGameModerations(s.T().Context(), gameID, publisher)

	s.Require().Error(err)
	s.Require().Equal(getGameErr, err)
}

func (s *TestSuite) TestGetGameModerations_GetPublisherError() {
	gameID := td.Int31()
	publisher := td.String()
	game := model.Game{ID: gameID, PublishersIDs: []int32{td.Int31()}}
	getPublisherErr := errors.New("get publisher error")

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(gomock.Any(), publisher).Return(int32(0), getPublisherErr)

	_, err := s.provider.GetGameModerations(s.T().Context(), gameID, publisher)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "get publisher id")
}

func (s *TestSuite) TestGetGameModerations_ForbiddenError() {
	gameID := td.Int31()
	publisher := td.String()
	publisherID := td.Int31()
	differentPublisherID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{differentPublisherID}, // Different publisher
	}

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(gomock.Any(), publisher).Return(publisherID, nil)

	_, err := s.provider.GetGameModerations(s.T().Context(), gameID, publisher)

	s.Require().Error(err)
	s.Require().True(apperr.IsStatusCode(err, apperr.Forbidden))
}

func (s *TestSuite) TestGetGameModerations_GetModerationsError() {
	gameID := td.Int31()
	publisher := td.String()
	publisherID := td.Int31()
	game := model.Game{
		ID:            gameID,
		PublishersIDs: []int32{publisherID},
	}
	getModerationsErr := errors.New("get moderations error")

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameID).Return(game, nil)
	s.storageMock.EXPECT().GetCompanyIDByName(gomock.Any(), publisher).Return(publisherID, nil)
	s.storageMock.EXPECT().GetModerationRecordsByGameID(gomock.Any(), gameID).Return(nil, getModerationsErr)

	_, err := s.provider.GetGameModerations(s.T().Context(), gameID, publisher)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "get moderations by game")
}
