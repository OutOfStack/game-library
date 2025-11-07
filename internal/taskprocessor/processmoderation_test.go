package taskprocessor_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestStartProcessModeration_Success() {
	task := model.Task{
		Name:     "process_moderation",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastProcessedGameId":0}`),
	}

	records := []model.ModerationIDGameID{
		{ModerationID: td.Int31(), GameID: td.Int31()},
		{ModerationID: td.Int31(), GameID: td.Int31()},
		{ModerationID: td.Int31(), GameID: td.Int31()},
	}

	moderationIDs := make([]int32, 0, len(records))
	for _, r := range records {
		moderationIDs = append(moderationIDs, r.ModerationID)
	}

	processErr := errors.New("process moderation failed")

	s.storageMock.EXPECT().
		RunWithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).
		Times(2)
	s.storageMock.EXPECT().GetTask(gomock.Any(), task.Name).Return(task, nil)

	firstUpdate := s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetPendingModerationGameIDs(gomock.Any(), 10).Return(records, nil)
	setInProgress := s.storageMock.EXPECT().SetModerationRecordsStatus(gomock.Any(), moderationIDs, model.ModerationStatusInProgress).Return(nil)

	call1 := s.moderationFacadeMock.EXPECT().ProcessModeration(gomock.Any(), records[0].GameID).After(setInProgress).Return(nil)
	call2 := s.moderationFacadeMock.EXPECT().ProcessModeration(gomock.Any(), records[1].GameID).After(setInProgress).Return(processErr)
	call3 := s.moderationFacadeMock.EXPECT().ProcessModeration(gomock.Any(), records[2].GameID).After(setInProgress).Return(nil)
	gomock.InOrder(call1, call2, call3)

	failedIDs := []int32{records[1].ModerationID}
	s.storageMock.EXPECT().
		SetModerationRecordsStatus(gomock.Any(), failedIDs, model.ModerationStatusPending).
		After(call3).
		Return(nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, updatedTask model.Task) error {
		s.Require().Equal(model.IdleTaskStatus, updatedTask.Status)

		var settings struct {
			LastProcessedGameID int32 `json:"lastProcessedGameId"`
		}
		s.Require().NoError(json.Unmarshal(updatedTask.Settings, &settings))

		var expectedLast int32
		for i, r := range records {
			if i == 1 {
				continue
			}
			if r.GameID > expectedLast {
				expectedLast = r.GameID
			}
		}
		s.Require().Equal(expectedLast, settings.LastProcessedGameID)

		return nil
	}).After(firstUpdate)

	err := s.provider.StartProcessModeration()
	s.Require().NoError(err)
}

func (s *TestSuite) TestStartProcessModeration_NoPendingRecords() {
	lastProcessed := td.Int31()
	task := model.Task{
		Name:     "process_moderation",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedGameId":%d}`, lastProcessed)),
	}

	s.storageMock.EXPECT().
		RunWithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).
		Times(2)
	s.storageMock.EXPECT().GetTask(gomock.Any(), task.Name).Return(task, nil)

	firstUpdate := s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetPendingModerationGameIDs(gomock.Any(), 10).Return([]model.ModerationIDGameID{}, nil)
	s.storageMock.EXPECT().SetModerationRecordsStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	s.moderationFacadeMock.EXPECT().ProcessModeration(gomock.Any(), gomock.Any()).Times(0)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, updatedTask model.Task) error {
		s.Require().Equal(model.IdleTaskStatus, updatedTask.Status)

		var settings struct {
			LastProcessedGameID int32 `json:"lastProcessedGameId"`
		}
		s.Require().NoError(json.Unmarshal(updatedTask.Settings, &settings))
		s.Require().Equal(lastProcessed, settings.LastProcessedGameID)

		return nil
	}).After(firstUpdate)

	err := s.provider.StartProcessModeration()
	s.Require().NoError(err)
}

func (s *TestSuite) TestStartProcessModeration_GetPendingError() {
	lastProcessed := td.Int31()
	task := model.Task{
		Name:     "process_moderation",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastProcessedGameId":%d}`, lastProcessed)),
	}

	expectedErr := errors.New("db failure")

	s.storageMock.EXPECT().
		RunWithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).
		Times(2)
	s.storageMock.EXPECT().GetTask(gomock.Any(), task.Name).Return(task, nil)

	firstUpdate := s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetPendingModerationGameIDs(gomock.Any(), 10).Return(nil, expectedErr)
	s.storageMock.EXPECT().SetModerationRecordsStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	s.moderationFacadeMock.EXPECT().ProcessModeration(gomock.Any(), gomock.Any()).Times(0)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, updatedTask model.Task) error {
		s.Require().Equal(model.ErrorTaskStatus, updatedTask.Status)
		s.Require().Equal(task.Settings, updatedTask.Settings)
		return nil
	}).After(firstUpdate)

	err := s.provider.StartProcessModeration()
	s.Require().NoError(err)
}
