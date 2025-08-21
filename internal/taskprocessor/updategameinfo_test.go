package taskprocessor_test

import (
	"errors"
	"fmt"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestStartUpdateGameInfo_Success() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31(), td.Int31()}
	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	game1 := model.Game{
		ID:     gameIDs[0],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}
	game2 := model.Game{
		ID:     gameIDs[1],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	updatedInfo1 := igdbapi.GameInfoForUpdate{
		ID:               game1.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID, platforms[1].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	updatedInfo2 := igdbapi.GameInfoForUpdate{
		ID:               game2.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 200).Return(gameIDs, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[0]).Return(game1, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game1.IGDBID).Return(updatedInfo1, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[0], gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[1]).Return(game2, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game2.IGDBID).Return(updatedInfo2, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[1], gomock.Any()).Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_NoGames() {
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastProcessedId":100}`),
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), int32(100), 200).Return([]int32{}, nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_GetGamesError() {
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastProcessedId":100}`),
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), int32(100), 200).Return(nil, errors.New("database error"))

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_GetPlatformsError() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31()}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 200).Return(gameIDs, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(nil, errors.New("platforms error"))

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_GetGameError() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31(), td.Int31()}
	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	game2 := model.Game{
		ID:     gameIDs[1],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	updatedInfo2 := igdbapi.GameInfoForUpdate{
		ID:               game2.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 200).Return(gameIDs, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[0]).Return(model.Game{}, errors.New("get game error"))
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[1]).Return(game2, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game2.IGDBID).Return(updatedInfo2, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[1], gomock.Any()).Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_SkipGameWithZeroIGDBID() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31(), td.Int31()}
	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	game1 := model.Game{
		ID:     gameIDs[0],
		Name:   td.String(),
		IGDBID: 0, // No IGDB ID, should be skipped
	}
	game2 := model.Game{
		ID:     gameIDs[1],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	updatedInfo2 := igdbapi.GameInfoForUpdate{
		ID:               game2.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 200).Return(gameIDs, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[0]).Return(game1, nil)
	// game1 should be skipped, no IGDB call for it
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[1]).Return(game2, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game2.IGDBID).Return(updatedInfo2, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[1], gomock.Any()).Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_IGDBAPIError() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31(), td.Int31()}
	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	game1 := model.Game{
		ID:     gameIDs[0],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}
	game2 := model.Game{
		ID:     gameIDs[1],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	updatedInfo2 := igdbapi.GameInfoForUpdate{
		ID:               game2.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 200).Return(gameIDs, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[0]).Return(game1, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game1.IGDBID).Return(igdbapi.GameInfoForUpdate{}, errors.New("igdb api error"))
	// game1 should continue to game2 despite error
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[1]).Return(game2, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game2.IGDBID).Return(updatedInfo2, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[1], gomock.Any()).Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_UpdateGameIGDBInfoError() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31(), td.Int31()}
	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	game1 := model.Game{
		ID:     gameIDs[0],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}
	game2 := model.Game{
		ID:     gameIDs[1],
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	updatedInfo1 := igdbapi.GameInfoForUpdate{
		ID:               game1.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	updatedInfo2 := igdbapi.GameInfoForUpdate{
		ID:               game2.IGDBID,
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int31(),
		Platforms:        []int64{platforms[0].IGDBID},
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
		},
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 200).Return(gameIDs, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)

	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[0]).Return(game1, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game1.IGDBID).Return(updatedInfo1, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[0], gomock.Any()).Return(errors.New("update error"))
	// continue to game2 despite error
	s.storageMock.EXPECT().GetGameByID(gomock.Any(), gameIDs[1]).Return(game2, nil)
	s.igdbClientMock.EXPECT().GetGameInfoForUpdate(gomock.Any(), game2.IGDBID).Return(updatedInfo2, nil)
	s.storageMock.EXPECT().UpdateGameIGDBInfo(gomock.Any(), gameIDs[1], gomock.Any()).Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateGameInfo_EmptySettings() {
	task := model.Task{
		Name:     "update_game_info",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte("{}"), // Empty JSON settings instead of nil
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	// Should start from ID 0 when settings is empty, return no games to process
	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), int32(0), 200).Return([]int32{}, nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateGameInfo()

	s.Require().NoError(err)
}
