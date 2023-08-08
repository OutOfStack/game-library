package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"github.com/OutOfStack/game-library/internal/client/uploadcare"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	minLengthForSearch = 2
)

// Game has handler methods for dealing with games
type Game struct {
	log        *zap.Logger
	storage    *repo.Storage
	igdb       *igdb.Client
	uploadcare *uploadcare.Client
	cache      *cache.Cache
}

// NewGame creates new Game
func NewGame(log *zap.Logger, storage *repo.Storage, igdb *igdb.Client, uploadcare *uploadcare.Client, cache *cache.Cache) *Game {
	return &Game{
		log:        log,
		storage:    storage,
		igdb:       igdb,
		uploadcare: uploadcare,
		cache:      cache,
	}
}

var tracer = otel.Tracer("")

var (
	companiesMap = cache.NewKVMap[int32, Company](1 * time.Hour)
	genresMap    = cache.NewKVMap[int32, Genre](1 * time.Hour)
	platformsMap = cache.NewKVMap[int32, Platform](0)
)

// GetGames godoc
// @Summary Get games
// @Description returns paginated games
// @ID get-games
// @Produce json
// @Param pageSize query int32  false "page size"
// @Param page     query int32  false "page"
// @Param orderBy  query string false "order by"	Enums(default, name, releaseDate)
// @Param name 	   query string false "name filter"
// @Success 200 {array}  GameResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [get]
func (g *Game) GetGames(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getgames")
	defer span.End()

	pageSizeParam := c.DefaultQuery("pageSize", "20")
	pageParam := c.DefaultQuery("page", "1")
	orderByParam := c.DefaultQuery("orderBy", "default")
	nameParam := c.DefaultQuery("name", "")

	pageSize, err := strconv.ParseInt(pageSizeParam, 10, 32)
	if err != nil || pageSize <= 0 {
		web.Err(c, web.NewRequestError(errors.New("Incorrect page size. Should be greater than 0"), http.StatusBadRequest))
		return
	}
	page, err := strconv.ParseInt(pageParam, 10, 32)
	if err != nil || page <= 0 {
		web.Err(c, web.NewRequestError(errors.New("Incorrect page. Should be greater than 0"), http.StatusBadRequest))
		return
	}
	var orderBy repo.OrderGamesBy
	switch orderByParam {
	case "default":
		orderBy = repo.OrderGamesByDefault
	case "name":
		orderBy = repo.OrderGamesByName
	case "releaseDate":
		orderBy = repo.OrderGamesByReleaseDate
	default:
		web.Err(c, web.NewRequestError(errors.New("Incorrect orderBy. Should be one of: default, releaseDate, name"), http.StatusBadRequest))
		return
	}
	name := nameParam
	if len(name) < minLengthForSearch {
		name = ""
	}

	span.SetAttributes(att.Int64("data.pageSize", pageSize), att.Int64("data.page", page),
		att.String("data.orderBy", orderByParam), att.String("data.name", nameParam))

	list := make([]repo.Game, 0)
	err = cache.Get(ctx, g.cache, getGamesKey(pageSize, page, string(orderBy), name), &list, func() ([]repo.Game, error) {
		return g.storage.GetGames(ctx, int(pageSize), int(page), orderBy, name)
	}, 0)
	if err != nil {
		web.Err(c, errors.Wrap(err, "getting games list"))
		return
	}

	response := make([]GameResponse, 0, len(list))
	for _, game := range list {
		r, err := g.mapToGameResponse(c, game)
		if err != nil {
			web.Err(c, web.NewRequestError(fmt.Errorf("error converting response"), http.StatusInternalServerError))
			return
		}
		response = append(response, r)
	}

	web.Respond(c, response, http.StatusOK)
}

// GetGamesCount godoc
// @Summary Get games count
// @Description returns games count
// @ID get-games-count
// @Produce json
// @Param name query string false "name filter"
// @Success 200 {array}  CountResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/count [get]
func (g *Game) GetGamesCount(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getgamescount")
	defer span.End()

	nameParam := c.DefaultQuery("name", "")
	name := nameParam
	if len(name) < minLengthForSearch {
		name = ""
	}
	span.SetAttributes(att.String("data.query", nameParam))

	var count uint64
	err := cache.Get(ctx, g.cache, getGamesCountKey(name), &count, func() (uint64, error) {
		return g.storage.GetGamesCount(ctx, name)
	}, 0)
	if err != nil {
		web.Err(c, errors.Wrap(err, "getting games count"))
		return
	}

	web.Respond(c, CountResponse{Count: count}, http.StatusOK)
}

// GetGame godoc
// @Summary Get game
// @Description returns game by ID
// @ID get-game-by-id
// @Produce json
// @Param 	id  path int32 true "Game ID"
// @Success 200 {object} GameResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [get]
func (g *Game) GetGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getgame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	var game repo.Game
	err = cache.Get(ctx, g.cache, getGameKey(id), &game, func() (repo.Game, error) {
		return g.storage.GetGameByID(ctx, id)
	}, 0)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, errors.Wrapf(err, "retrieving game with id %v", id))
		return
	}

	resp, err := g.mapToGameResponse(c, game)
	if err != nil {
		web.Err(c, web.NewRequestError(fmt.Errorf("error converting response"), http.StatusInternalServerError))
		return
	}
	web.Respond(c, resp, http.StatusOK)
}

// CreateGame godoc
// @Summary Create game
// @Description creates new game
// @ID create-game
// @Accept  json
// @Produce json
// @Param   game body CreateGameRequest true "create game"
// @Success 201 {object} IDResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [post]
func (g *Game) CreateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.creategame")
	defer span.End()

	var cg CreateGameRequest
	err := web.Decode(c, &cg)
	if err != nil {
		web.Err(c, errors.Wrap(err, "decoding new game"))
		return
	}
	span.SetAttributes(att.String("data.name", cg.Name))

	claims, err := web.GetClaims(c)
	if err != nil {
		web.Err(c, errors.Wrap(err, "getting claims from context"))
		return
	}

	developer, publisher := cg.Developer, claims.Name
	// get id or create developer
	developerID, err := g.storage.GetCompanyIDByName(ctx, developer)
	if err != nil && !errors.As(err, &repo.ErrNotFound[string]{}) {
		web.Err(c, errors.Wrapf(err, "get company id with name %s", developer))
		return
	}
	if developerID == 0 {
		developerID, err = g.storage.CreateCompany(ctx, repo.Company{
			Name: developer,
		})
		if err != nil {
			web.Err(c, errors.Wrapf(err, "create company %s", developer))
			return
		}
	}

	// get id or create publisher
	publisherID, err := g.storage.GetCompanyIDByName(ctx, publisher)
	if err != nil && !errors.As(err, &repo.ErrNotFound[string]{}) {
		web.Err(c, errors.Wrapf(err, "get company id with name %s", publisher))
		return
	}
	if publisherID == 0 {
		publisherID, err = g.storage.CreateCompany(ctx, repo.Company{
			Name: publisher,
		})
		if err != nil {
			web.Err(c, errors.Wrapf(err, "create company %s", publisher))
			return
		}
	}

	create := mapToCreateGame(&cg, developerID, publisherID)

	id, err := g.storage.CreateGame(ctx, create)
	if err != nil {
		web.Err(c, errors.Wrap(err, "adding new game"))
		return
	}

	// invalidate cache
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate games count cache
		key = gamesCountKey
		err = cache.DeleteByStartsWith(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
	}()

	web.Respond(c, IDResponse{ID: id}, http.StatusCreated)
}

// UpdateGame godoc
// @Summary Update game
// @Description updates game by ID
// @ID update-game
// @Accept  json
// @Produce json
// @Param  	id   path int32 			true "Game ID"
// @Param  	game body UpdateGameRequest true "update game"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [patch]
func (g *Game) UpdateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.updategame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	var ugr UpdateGameRequest
	if err = web.Decode(c, &ugr); err != nil {
		web.Err(c, errors.Wrap(err, "decoding game update"))
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	game, err := g.storage.GetGameByID(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, errors.Wrapf(err, "retrieve game with id %d", id))
		return
	}

	developer := ugr.Developer
	developers := game.Developers
	if developer != nil {
		if *developer != "" {
			// get id or create developer
			developerID, err := g.storage.GetCompanyIDByName(ctx, *developer)
			if err != nil && !errors.As(err, &repo.ErrNotFound[string]{}) {
				web.Err(c, errors.Wrapf(err, "get developer id with name %s", *developer))
				return
			}
			if developerID == 0 {
				developerID, err = g.storage.CreateCompany(ctx, repo.Company{
					Name: *developer,
				})
				if err != nil {
					web.Err(c, errors.Wrapf(err, "create developer %s", *developer))
					return
				}
			}
			developers = []int32{developerID}
		} else {
			developers = []int32{}
		}
	}

	update := mapToUpdateGame(game, ugr)
	update.Developers = developers

	err = g.storage.UpdateGame(ctx, id, update)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, errors.Wrapf(err, "updating game with id %v", id))
		return
	}

	// invalidate cache
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}

		// invalidate game cache
		key = getGameKey(id)
		err = cache.Delete(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by key", zap.String("key", key), zap.Error(err))
		}
		// recache game
		err = cache.Get(ctx, g.cache, key, new(repo.Game), func() (repo.Game, error) {
			return g.storage.GetGameByID(ctx, id)
		}, 0)
		if err != nil {
			g.log.Error("recache game with id", zap.Int32("id", id), zap.Error(err))
		}
	}()

	web.Respond(c, nil, http.StatusNoContent)
}

// DeleteGame godoc
// @Summary Delete game
// @Description deletes game by ID
// @ID delete-game
// @Accept  json
// @Produce json
// @Param  	id path int32 true "Game ID"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [delete]
func (g *Game) DeleteGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.game.delete")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	err = g.storage.DeleteGame(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, errors.Wrapf(err, "deleting game with id %v", id))
		return
	}

	// invalidate cache
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate games count cache
		key = gamesCountKey
		err = cache.DeleteByStartsWith(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate game cache
		key = getGameKey(id)
		err = cache.Delete(ctx, g.cache, key)
		if err != nil {
			g.log.Error("remove game cache by key", zap.String("key", key), zap.Error(err))
		}
	}()

	web.Respond(c, nil, http.StatusNoContent)
}
