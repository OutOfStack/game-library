package taskprocessor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Storage db storage interface
type Storage interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	GetTask(ctx context.Context, tx *sqlx.Tx, name string) (repo.Task, error)
	UpdateTask(ctx context.Context, tx *sqlx.Tx, task repo.Task) error
}

// IGDBProvider igdb client interface
type IGDBProvider interface {
}

// TaskProvider contains dependencies for tasks
type TaskProvider struct {
	log          *zap.Logger
	storage      Storage
	igdbProvider IGDBProvider
}

// New creates new TaskProvider
func New(log *zap.Logger, storage Storage, igdbProvider IGDBProvider) *TaskProvider {
	return &TaskProvider{
		log:          log,
		storage:      storage,
		igdbProvider: igdbProvider,
	}
}

func (tp *TaskProvider) DoTask(name string, taskFn func() error) error {
	ctx := context.Background()

	tx, err := tp.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %v", err)
	}

	task, err := tp.storage.GetTask(ctx, tx, name)
	if err != nil {
		if errors.Is(err, repo.ErrTransactionLocked) {
			return tx.Rollback()
		}
		return err
	}

	if task.Status == repo.RunningTaskStatus {
		return tx.Rollback()
	}

	task.Status = repo.RunningTaskStatus
	task.RunCount++
	task.LastRun = sql.NullTime{Time: time.Now(), Valid: true}

	err = tp.storage.UpdateTask(ctx, tx, task)
	if err != nil {
		tp.log.Error("update task", zap.Error(err))
		rErr := tx.Rollback()
		if rErr != nil {
			return rErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = taskFn()
	if err != nil {
		tp.log.Error("running task", zap.Error(err))
		task.Status = repo.ErrorTaskStatus
	} else {
		task.Status = repo.IdleTaskStatus
	}

	return tp.storage.UpdateTask(ctx, nil, task)
}
