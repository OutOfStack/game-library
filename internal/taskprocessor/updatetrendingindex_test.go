package taskprocessor_test

import (
	"errors"
	"fmt"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestStartUpdateTrendingIndex_Success() {
	lastProcessedID := td.Int31()
	task := model.Task{
		Name:     "update_trending_index",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedId":%d}`, lastProcessedID)),
	}

	gameIDs := []int32{td.Int31(), td.Int31(), td.Int31()}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), lastProcessedID, 300).Return(gameIDs, nil)
	for _, gameID := range gameIDs {
		s.gameFacadeMock.EXPECT().UpdateGameTrendingIndex(gomock.Any(), gameID).Return(nil)
	}

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateTrendingIndex()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateTrendingIndex_NoGames() {
	task := model.Task{
		Name:     "update_trending_index",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastProcessedId":100}`),
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), int32(100), 300).Return([]int32{}, nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateTrendingIndex()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateTrendingIndex_GetGamesError() {
	task := model.Task{
		Name:     "update_trending_index",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastProcessedId":100}`),
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), int32(100), 300).Return(nil, errors.New("database error"))

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateTrendingIndex()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartUpdateTrendingIndex_UpdateGameError() {
	task := model.Task{
		Name:     "update_trending_index",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastProcessedId":100}`),
	}

	gameIDs := []int32{td.Int31(), td.Int31()}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetGamesIDsAfterID(gomock.Any(), int32(100), 300).Return(gameIDs, nil)

	s.gameFacadeMock.EXPECT().UpdateGameTrendingIndex(gomock.Any(), gameIDs[0]).Return(errors.New("update error"))
	s.gameFacadeMock.EXPECT().UpdateGameTrendingIndex(gomock.Any(), gameIDs[1]).Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartUpdateTrendingIndex()

	s.Require().NoError(err)
}
