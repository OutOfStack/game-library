package taskprocessor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Storage db storage interface
type Storage interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	GetTask(ctx context.Context, tx *sqlx.Tx, name string) (repo.Task, error)
	UpdateTask(ctx context.Context, tx *sqlx.Tx, task repo.Task) error
	CreateGame(ctx context.Context, cg repo.CreateGame) (id int32, err error)
	GetGameIDByIGDBID(ctx context.Context, igdbID int64) (id int32, err error)
	GetPlatforms(ctx context.Context) ([]repo.Platform, error)
	CreateGenre(ctx context.Context, g repo.Genre) (int32, error)
	GetGenres(ctx context.Context) ([]repo.Genre, error)
	CreateCompany(ctx context.Context, c repo.Company) (int32, error)
	GetCompanies(ctx context.Context) ([]repo.Company, error)
}

// IGDBClient igdb client interface
type IGDBClient interface {
	GetTopRatedGames(ctx context.Context, platformsIDs []int64, releasedAfter time.Time, minRatingsCount, minRating, limit int64) ([]igdb.TopRatedGamesResp, error)
}

// UploadcareClient uploadcare client interface
type UploadcareClient interface {
	UploadImageFromURL(ctx context.Context, imageURL string) (newURL string, err error)
}

// TaskProvider contains dependencies for tasks
type TaskProvider struct {
	log                *zap.Logger
	storage            Storage
	igdbProvider       IGDBClient
	uploadcareProvider UploadcareClient
}

// New creates new TaskProvider
func New(log *zap.Logger, storage Storage, igdbClient IGDBClient, uploadcareClient UploadcareClient) *TaskProvider {
	return &TaskProvider{
		log:                log,
		storage:            storage,
		igdbProvider:       igdbClient,
		uploadcareProvider: uploadcareClient,
	}
}

// DoTask - runs task.
// name - name of a task to run.
// taskFn - function to be run: it accepts settings and returns updates settings
func (tp *TaskProvider) DoTask(name string, taskFn func(ctx context.Context, settings repo.TaskSettings) (newSettings repo.TaskSettings, err error)) error {
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

	settings := make([]byte, len(task.Settings))
	copy(settings, task.Settings)

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

	tp.log.Info("task started", zap.String("name", name))
	task.Settings, err = taskFn(ctx, settings)
	if err != nil {
		tp.log.Error("running task", zap.Error(err))
		task.Status = repo.ErrorTaskStatus
	} else {
		task.Status = repo.IdleTaskStatus
	}
	tp.log.Info("task finished", zap.String("name", name), zap.Error(err))

	return tp.storage.UpdateTask(ctx, nil, task)
}
