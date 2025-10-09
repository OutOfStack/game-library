package taskprocessor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const (
	taskTimeout = 15 * time.Minute

	// igdb api rate limit (rps)
	igdbAPIRPSLimit = 4
)

// Storage db storage interface
type Storage interface {
	RunWithTx(ctx context.Context, f func(context.Context) error) error

	GetTask(ctx context.Context, name string) (model.Task, error)
	UpdateTask(ctx context.Context, task model.Task) error

	CreateGame(ctx context.Context, cgd model.CreateGameData) (id int32, err error)
	GetGameIDByIGDBID(ctx context.Context, igdbID int64) (id int32, err error)
	GetGameByID(ctx context.Context, id int32) (game model.Game, err error)
	UpdateGameIGDBInfo(ctx context.Context, id int32, ug model.UpdateGameIGDBData) error
	GetPlatforms(ctx context.Context) ([]model.Platform, error)
	CreateGenre(ctx context.Context, g model.Genre) (int32, error)
	GetGenres(ctx context.Context) ([]model.Genre, error)
	CreateCompany(ctx context.Context, c model.Company) (int32, error)
	GetCompanies(ctx context.Context) ([]model.Company, error)
	GetGamesIDsAfterID(ctx context.Context, lastID int32, batchSize int) ([]int32, error)

	GetPendingModerationGameIDs(ctx context.Context, limit int) ([]model.ModerationIDGameID, error)
	SetModerationRecordsStatus(ctx context.Context, gameIDs []int32, status model.ModerationStatus) error
}

// IGDBAPIClient igdb api client interface
type IGDBAPIClient interface {
	GetTopRatedGames(ctx context.Context, platformsIDs []int64, releasedAfter time.Time, minRatingsCount, minRating, limit int64) ([]igdbapi.TopRatedGames, error)
	GetImageByURL(ctx context.Context, imageURL, imageType string) (igdbapi.GetImageResp, error)
	GetGameInfoForUpdate(ctx context.Context, igdbID int64) (igdbapi.GameInfoForUpdate, error)
}

// S3Client s3 store client interface
type S3Client interface {
	Upload(ctx context.Context, data io.ReadSeeker, contentType string, md map[string]string) (s3.UploadResult, error)
}

// GameFacade game facade interface
type GameFacade interface {
	UpdateGameTrendingIndex(ctx context.Context, gameID int32) error
}

// ModerationFacade moderation facade interface
type ModerationFacade interface {
	ProcessModeration(ctx context.Context, gameID int32) error
}

// TaskProvider contains dependencies for tasks
type TaskProvider struct {
	log              *zap.Logger
	storage          Storage
	igdbAPIClient    IGDBAPIClient
	s3Client         S3Client
	gameFacade       GameFacade
	moderationFacade ModerationFacade
	igdbAPILimiter   *rate.Limiter
}

// New creates new TaskProvider
func New(log *zap.Logger, storage Storage, igdbClient IGDBAPIClient, s3Client S3Client, gameFacade GameFacade, moderationFacade ModerationFacade) *TaskProvider {
	return &TaskProvider{
		log:              log,
		storage:          storage,
		igdbAPIClient:    igdbClient,
		igdbAPILimiter:   rate.NewLimiter(rate.Every(time.Second), igdbAPIRPSLimit),
		s3Client:         s3Client,
		gameFacade:       gameFacade,
		moderationFacade: moderationFacade,
	}
}

// DoTask - runs task.
// name - name of a task to run.
// taskFn - function to be run: it accepts settings and returns updates settings
func (tp *TaskProvider) DoTask(name string, taskFn func(ctx context.Context, settings model.TaskSettings) (newSettings model.TaskSettings, err error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), taskTimeout)
	defer cancel()

	var settings []byte
	var task model.Task
	var err error

	txErr := tp.storage.RunWithTx(ctx, func(ctx context.Context) error {
		task, err = tp.storage.GetTask(ctx, name)
		if err != nil {
			tp.log.Error("get task", zap.String("name", name), zap.Error(err))
			return err
		}

		if task.Status == model.RunningTaskStatus {
			return fmt.Errorf("task %s is already running", name)
		}

		task.Status = model.RunningTaskStatus
		task.RunCount++
		task.LastRun = sql.NullTime{Time: time.Now(), Valid: true}

		settings = make([]byte, len(task.Settings))
		copy(settings, task.Settings)

		err = tp.storage.UpdateTask(ctx, task)
		if err != nil {
			tp.log.Error("update task", zap.String("name", name), zap.Error(err))
			return err
		}

		return nil
	})
	if txErr != nil {
		if errors.Is(txErr, repo.ErrTransactionLocked) {
			return nil
		}
		return txErr
	}

	tp.log.Info("task started", zap.String("name", name))

	task.Settings, err = taskFn(ctx, settings)
	if err != nil {
		tp.log.Error("run task", zap.Error(err))
		task.Status = model.ErrorTaskStatus
	} else {
		task.Status = model.IdleTaskStatus
	}

	tp.log.Info("task finished", zap.String("name", name), zap.Error(err))

	return tp.storage.UpdateTask(ctx, task)
}
