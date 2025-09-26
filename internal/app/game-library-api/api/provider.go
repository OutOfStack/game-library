package api

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("api")

// GameFacade represents methods for working with games
type GameFacade interface {
	GetGames(ctx context.Context, page, pageSize uint32, filter model.GamesFilter) (games []model.Game, count uint64, err error)
	GetGameByID(ctx context.Context, id int32) (model.Game, error)
	CreateGame(ctx context.Context, cg model.CreateGame) (id int32, err error)
	UpdateGame(ctx context.Context, id int32, upd model.UpdateGame) error
	DeleteGame(ctx context.Context, id int32, publisher string) error
	RateGame(ctx context.Context, gameID int32, userID string, rating uint8) error
	GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error)
	UploadGameImages(ctx context.Context, coverFiles, screenshotFiles []*multipart.FileHeader, publisherName string) ([]model.File, error)

	GetGenres(ctx context.Context) ([]model.Genre, error)
	GetGenresMap(ctx context.Context) (map[int32]model.Genre, error)
	GetTopGenres(ctx context.Context, limit int64) ([]model.Genre, error)

	GetPlatforms(ctx context.Context) ([]model.Platform, error)
	GetPlatformsMap(ctx context.Context) (map[int32]model.Platform, error)

	GetCompaniesMap(ctx context.Context) (map[int32]model.Company, error)
	GetTopCompanies(ctx context.Context, companyType string, limit int64) ([]model.Company, error)

	GetPublisherGames(ctx context.Context, publisher string) ([]model.Game, error)
	GetGameModerations(ctx context.Context, gameID int32, publisher string) ([]model.Moderation, error)
}

// Decoder decodes request
type Decoder interface {
	Decode(r *http.Request, val interface{}) error
}

// Provider has all dependencies for handlers
type Provider struct {
	log        *zap.Logger
	cache      *cache.RedisStore
	gameFacade GameFacade
	decoder    Decoder
}

// NewProvider creates new provider
func NewProvider(log *zap.Logger, redisStore *cache.RedisStore, gameFacade GameFacade, decoder Decoder) *Provider {
	return &Provider{
		log:        log,
		cache:      redisStore,
		gameFacade: gameFacade,
		decoder:    decoder,
	}
}
