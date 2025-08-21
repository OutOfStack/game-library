package facade

import (
	"context"
	"io"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"go.uber.org/zap"
)

// Provider represents dependencies for facade layer
type Provider struct {
	log      *zap.Logger
	storage  Storage
	cache    *cache.RedisStore
	s3Client S3Client
}

// NewProvider returns nre provider
func NewProvider(logger *zap.Logger, storage Storage, cache *cache.RedisStore, s3Client S3Client) *Provider {
	return &Provider{
		log:      logger,
		storage:  storage,
		cache:    cache,
		s3Client: s3Client,
	}
}

// Storage provides methods for working with database
type Storage interface {
	GetGames(ctx context.Context, pageSize, page uint32, filter model.GamesFilter) (list []model.Game, err error)
	GetGamesCount(ctx context.Context, filter model.GamesFilter) (count uint64, err error)
	GetGameByID(ctx context.Context, id int32) (game model.Game, err error)
	CreateGame(ctx context.Context, cg model.CreateGameData) (id int32, err error)
	UpdateGame(ctx context.Context, id int32, ug model.UpdateGameData) error
	DeleteGame(ctx context.Context, id int32) error
	UpdateGameRating(ctx context.Context, id int32) error
	GetPublisherGamesCount(ctx context.Context, publisherID int32, startDate, endDate time.Time) (count int, err error)
	UpdateGameTrendingIndex(ctx context.Context, gameID int32, trendingIndex float64) error
	GetGameTrendingData(ctx context.Context, gameID int32) (model.GameTrendingData, error)

	CreateCompany(ctx context.Context, c model.Company) (id int32, err error)
	GetCompanies(ctx context.Context) (companies []model.Company, err error)
	GetCompanyByID(ctx context.Context, id int32) (company model.Company, err error)
	GetCompanyIDByName(ctx context.Context, name string) (id int32, err error)
	GetTopDevelopers(ctx context.Context, limit int64) (companies []model.Company, err error)
	GetTopPublishers(ctx context.Context, limit int64) (companies []model.Company, err error)

	GetGenres(ctx context.Context) (genres []model.Genre, err error)
	GetGenreByID(ctx context.Context, id int32) (genre model.Genre, err error)
	GetTopGenres(ctx context.Context, limit int64) (genres []model.Genre, err error)

	GetPlatforms(ctx context.Context) (platforms []model.Platform, err error)
	GetPlatformByID(ctx context.Context, id int32) (platform model.Platform, err error)

	AddRating(ctx context.Context, cr model.CreateRating) error
	RemoveRating(ctx context.Context, rr model.RemoveRating) error
	GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error)
}

// S3Client represents the interface for S3 client operations
type S3Client interface {
	Upload(ctx context.Context, data io.ReadSeeker, contentType string, md map[string]string) (s3.UploadResult, error)
}
