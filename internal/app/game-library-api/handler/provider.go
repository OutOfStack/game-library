package handler

import (
	"context"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("")

// Provider has all dependencies for handlers
type Provider struct {
	log     *zap.Logger
	storage Storage
	cache   *cache.RedisStore
}

// NewProvider creates new provider
func NewProvider(log *zap.Logger, storage Storage, redisStore *cache.RedisStore) *Provider {
	return &Provider{
		log:     log,
		storage: storage,
		cache:   redisStore,
	}
}

var (
	companiesMap = cache.NewKVMap[int32, Company](1 * time.Hour)
	genresMap    = cache.NewKVMap[int32, Genre](1 * time.Hour)
	platformsMap = cache.NewKVMap[int32, Platform](0)
)

// Storage provides methods for working with repo
type Storage interface {
	GetGames(ctx context.Context, pageSize, page int, filter repo.GamesFilter) (list []repo.Game, err error)
	GetGamesCount(ctx context.Context, filter repo.GamesFilter) (count uint64, err error)
	GetGameByID(ctx context.Context, id int32) (game repo.Game, err error)
	CreateGame(ctx context.Context, cg repo.CreateGame) (id int32, err error)
	UpdateGame(ctx context.Context, id int32, ug repo.UpdateGame) error
	DeleteGame(ctx context.Context, id int32) error
	UpdateGameRating(ctx context.Context, id int32) error

	CreateCompany(ctx context.Context, c repo.Company) (id int32, err error)
	GetCompanies(ctx context.Context) (companies []repo.Company, err error)
	GetCompanyIDByName(ctx context.Context, name string) (id int32, err error)
	GetTopDevelopers(ctx context.Context, limit int64) (companies []repo.Company, err error)
	GetTopPublishers(ctx context.Context, limit int64) (companies []repo.Company, err error)

	GetGenres(ctx context.Context) (genres []repo.Genre, err error)

	GetPlatforms(ctx context.Context) (platforms []repo.Platform, err error)

	AddRating(ctx context.Context, cr repo.CreateRating) error
	RemoveRating(ctx context.Context, rr repo.RemoveRating) error
	GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error)
}
