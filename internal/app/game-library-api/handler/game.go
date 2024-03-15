package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/gin-gonic/gin"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	minLengthForSearch = 2
)

// GetGames godoc
// @Summary Get games
// @Description returns paginated games
// @ID get-games
// @Produce json
// @Param pageSize  query int32  false "page size"
// @Param page      query int32  false "page"
// @Param orderBy   query string false "order by"	Enums(default, name, releaseDate)
// @Param name 	    query string false "name filter"
// @Param genre     query int32  false "genre filter"
// @Param developer query int32  false "developer id filter"
// @Param publisher query int32  false "publisher id filter"
// @Success 200 {object}  GamesResponse
// @Failure 500 {object}  web.ErrorResponse
// @Router /games [get]
func (p *Provider) GetGames(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getGames")
	defer span.End()

	for key, values := range c.Request.URL.Query() {
		if len(values) == 1 {
			span.SetAttributes(att.String(key, values[0]))
		} else if len(values) > 1 {
			span.SetAttributes(att.StringSlice(key, values))
		}
	}

	// get query params and form filter
	var queryParams GetGamesQueryParams
	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		web.Err(c, web.NewRequestError(errors.New("incorrect query params"), http.StatusBadRequest))
		return
	}

	page, pageSize := queryParams.Page, queryParams.PageSize
	var filter repo.GamesFilter
	if len(queryParams.Name) >= minLengthForSearch {
		filter.Name = queryParams.Name
	}
	if queryParams.Genre != 0 {
		filter.GenreID = queryParams.Genre
	}
	if queryParams.Developer != 0 {
		filter.DeveloperID = queryParams.Developer
	}
	if queryParams.Publisher != 0 {
		filter.PublisherID = queryParams.Publisher
	}
	switch queryParams.OrderBy {
	case "", "default":
		filter.OrderBy = repo.OrderGamesByDefault
	case "name":
		filter.OrderBy = repo.OrderGamesByName
	case "releaseDate":
		filter.OrderBy = repo.OrderGamesByReleaseDate
	default:
		web.Err(c, web.NewRequestError(errors.New("incorrect orderBy. Should be one of: default, releaseDate, name"), http.StatusBadRequest))
		return
	}

	list := make([]repo.Game, 0)
	err = cache.Get(ctx, p.cache, getGamesKey(int64(pageSize), int64(page), filter), &list, func() ([]repo.Game, error) {
		return p.storage.GetGames(ctx, pageSize, page, filter)
	}, 0)
	if err != nil {
		web.Err(c, fmt.Errorf("getting games list: %w", err))
		return
	}

	var count uint64
	err = cache.Get(ctx, p.cache, getGamesCountKey(filter), &count, func() (uint64, error) {
		return p.storage.GetGamesCount(ctx, filter)
	}, 0)
	if err != nil {
		web.Err(c, fmt.Errorf("getting games count: %w", err))
		return
	}

	games := make([]GameResponse, 0, len(list))
	for _, game := range list {
		r, mErr := p.mapToGameResponse(c, game)
		if mErr != nil {
			web.Err(c, web.NewRequestError(fmt.Errorf("error converting response"), http.StatusInternalServerError))
			return
		}
		games = append(games, r)
	}

	response := GamesResponse{
		Games: games,
		Count: count,
	}
	web.Respond(c, response, http.StatusOK)
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
func (p *Provider) GetGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getGame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	var game repo.Game
	err = cache.Get(ctx, p.cache, getGameKey(id), &game, func() (repo.Game, error) {
		return p.storage.GetGameByID(ctx, id)
	}, 0)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, fmt.Errorf("retrieving game with id %d: %w", id, err))
		return
	}

	resp, err := p.mapToGameResponse(c, game)
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
func (p *Provider) CreateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.createGame")
	defer span.End()

	var cg CreateGameRequest
	err := web.Decode(c, &cg)
	if err != nil {
		web.Err(c, fmt.Errorf("decoding new game: %w", err))
		return
	}
	span.SetAttributes(att.String("data.name", cg.Name))

	claims, err := web.GetClaims(c)
	if err != nil {
		web.Err(c, fmt.Errorf("getting claims from context: %w", err))
		return
	}

	developer, publisher := cg.Developer, claims.Name
	// get developer id or create developer
	developerID, err := p.storage.GetCompanyIDByName(ctx, developer)
	if err != nil && !errors.As(err, &repo.ErrNotFound[string]{}) {
		web.Err(c, fmt.Errorf("get company id with name %s: %w", developer, err))
		return
	}
	if developerID == 0 {
		developerID, err = p.storage.CreateCompany(ctx, repo.Company{
			Name: developer,
		})
		if err != nil {
			web.Err(c, fmt.Errorf("create company %s: %w", developer, err))
			return
		}
	}

	// get id or create publisher
	publisherID, err := p.storage.GetCompanyIDByName(ctx, publisher)
	if err != nil && !errors.As(err, &repo.ErrNotFound[string]{}) {
		web.Err(c, fmt.Errorf("get company id with name %s: %w", publisher, err))
		return
	}
	if publisherID == 0 {
		publisherID, err = p.storage.CreateCompany(ctx, repo.Company{
			Name: publisher,
		})
		if err != nil {
			web.Err(c, fmt.Errorf("create company %s: %w", publisher, err))
			return
		}
	}

	create := mapToCreateGame(&cg, developerID, publisherID)

	id, err := p.storage.CreateGame(ctx, create)
	if err != nil {
		web.Err(c, fmt.Errorf("adding new game: %w", err))
		return
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 1*time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate games count cache
		key = gamesCountKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}

		// invalidate companies cache
		companiesMap.Purge()
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
func (p *Provider) UpdateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.updateGame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	var ugr UpdateGameRequest
	if err = web.Decode(c, &ugr); err != nil {
		web.Err(c, fmt.Errorf("decoding game update: %w", err))
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	game, err := p.storage.GetGameByID(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, fmt.Errorf("retrieve game with id %d: %w", id, err))
		return
	}

	developer := ugr.Developer
	developers := game.Developers
	if developer != nil {
		if *developer == "" {
			developers = []int32{}
		} else {
			// get id or create developer
			developerID, cErr := p.storage.GetCompanyIDByName(ctx, *developer)
			if cErr != nil && !errors.As(cErr, &repo.ErrNotFound[string]{}) {
				web.Err(c, fmt.Errorf("get developer id with name %s: %w", *developer, cErr))
				return
			}
			if developerID == 0 {
				developerID, err = p.storage.CreateCompany(ctx, repo.Company{
					Name: *developer,
				})
				if err != nil {
					web.Err(c, fmt.Errorf("create developer %s: %w", *developer, err))
					return
				}
			}
			developers = []int32{developerID}
		}
	}

	update := mapToUpdateGame(game, ugr)
	update.Developers = developers

	err = p.storage.UpdateGame(ctx, id, update)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, fmt.Errorf("updating game with id %v: %w", id, err))
		return
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 1*time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}

		// invalidate game cache
		key = getGameKey(id)
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by key", zap.String("key", key), zap.Error(err))
		}
		// recache game
		err = cache.Get(bCtx, p.cache, key, new(repo.Game), func() (repo.Game, error) {
			return p.storage.GetGameByID(bCtx, id)
		}, 0)
		if err != nil {
			p.log.Error("recache game with id", zap.Int32("id", id), zap.Error(err))
		}

		// invalidate companies cache
		companiesMap.Purge()
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
func (p *Provider) DeleteGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.deleteGame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	err = p.storage.DeleteGame(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, fmt.Errorf("deleting game with id %v: %v", id, err))
		return
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 1*time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate games count cache
		key = gamesCountKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate game cache
		key = getGameKey(id)
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove game cache by key", zap.String("key", key), zap.Error(err))
		}
	}()

	web.Respond(c, nil, http.StatusNoContent)
}
