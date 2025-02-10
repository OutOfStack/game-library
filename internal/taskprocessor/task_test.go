package taskprocessor_test

import (
	"context"
	"errors"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestDoTask_Success() {
	task := model.Task{
		Name:     td.String(),
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastReleasedAt":"2025-01-01T00:00:00Z"}`),
	}

	taskFn := func(_ context.Context, _ model.TaskSettings) (model.TaskSettings, error) {
		return []byte(`{"lastReleasedAt":"2025-01-02T00:00:00Z"}`), nil
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.DoTask(task.Name, taskFn)

	s.Require().NoError(err)
}

func (s *TestSuite) TestDoTask_ErrorOnBeginTx() {
	taskName := td.String()
	taskFn := func(_ context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
		return settings, nil
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(nil, errors.New("new error"))

	err := s.provider.DoTask(taskName, taskFn)

	s.Require().Error(err)
	s.Contains(err.Error(), "begin transaction")
}

func (s *TestSuite) TestDoTask_TransactionLocked() {
	taskName := td.String()
	taskFn := func(_ context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
		return settings, nil
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, taskName).Return(model.Task{}, repo.ErrTransactionLocked)
	s.tx.EXPECT().Rollback(gomock.Any()).Return(nil)

	err := s.provider.DoTask(taskName, taskFn)

	s.Require().NoError(err)
}

func (s *TestSuite) TestDoTask_ErrorOnUpdateTask() {
	task := model.Task{
		Name:     td.String(),
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(`{"lastReleasedAt":"2025-01-01T00:00:00Z"}`),
	}
	updateTaskErr := errors.New("update error")

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(updateTaskErr)
	s.tx.EXPECT().Rollback(gomock.Any()).Return(nil)

	taskFn := func(_ context.Context, settings model.TaskSettings) (model.TaskSettings, error) {
		return settings, nil
	}

	err := s.provider.DoTask(task.Name, taskFn)

	s.Require().ErrorIs(err, updateTaskErr)
}
