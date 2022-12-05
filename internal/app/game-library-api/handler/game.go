package handler

import (
	"net/http"
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// Game has handler methods for dealing with games
type Game struct {
	Log     *zap.Logger
	Storage *repo.Storage
	IGDB    *igdb.Client
}

var tracer = otel.Tracer("")

// GetGames godoc
// @Summary Get games
// @Description returns paginated games
// @ID get-games
// @Produce json
// @Param pageSize query int32 false "page size"
// @Param lastId   query int32 false "last fetched Id"
// @Success 200 {array}  GameResp
// @Failure 500 {object} web.ErrorResponse
// @Router /games [get]
func (g *Game) GetGames(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getgames")
	defer span.End()

	psParam := c.DefaultQuery("pageSize", "20")
	liParam := c.DefaultQuery("lastId", "0")
	pageSize, err := strconv.ParseInt(psParam, 10, 32)
	if err != nil || pageSize <= 0 {
		c.Error(web.NewRequestError(errors.New("Incorrect page size. Should be greater than 0"), http.StatusBadRequest))
		return
	}
	lastID, err := strconv.ParseInt(liParam, 10, 32)
	if err != nil || lastID < 0 {
		c.Error(web.NewRequestError(errors.New("Incorrect last Id. Should be greater or equal to 0"), http.StatusBadRequest))
		return
	}
	span.SetAttributes(attribute.Int64("data.pagesize", pageSize), attribute.Int64("data.lastid", lastID))

	list, err := g.Storage.GetGames(ctx, int(pageSize), int32(lastID))

	if err != nil {
		c.Error(errors.Wrap(err, "getting games list"))
		return
	}

	resps := []GameResp{}
	for _, g := range list {
		resps = append(resps, mapGameToResp(g))
	}

	web.Respond(c, resps, http.StatusOK)
}

// GetGame godoc
// @Summary Get game
// @Description returns game by ID
// @ID get-game-by-id
// @Produce json
// @Param 	id  path int32 true "Game ID"
// @Success 200 {object} GameResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [get]
func (g *Game) GetGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getgame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(id)))

	game, err := g.Storage.GetGameByID(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "retrieving game with id %v", id))
		return
	}

	resp := mapGameToResp(game)
	web.Respond(c, resp, http.StatusOK)
}

// SearchGames godoc
// @Summary Searches games by name
// @Description returns games filtered by provided name
// @ID search-games
// @Produce json
// @Param name query string false "name to search by"
// @Success 200 {array}  GameResp
// @Failure 500 {object} web.ErrorResponse
// @Router /games/search [get]
func (g *Game) SearchGames(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.searchgames")
	defer span.End()

	nameParam := c.DefaultQuery("name", "")
	if len(nameParam) < 2 {
		c.Error(web.NewRequestError(errors.New("Length of name to be searched by should be at least 2 characters"), http.StatusBadRequest))
		return
	}
	span.SetAttributes(attribute.String("data.query", nameParam))

	list, err := g.Storage.SearchGames(ctx, nameParam)
	if err != nil {
		c.Error(errors.Wrap(err, "searching games list"))
		return
	}

	resps := []GameResp{}
	for _, g := range list {
		resps = append(resps, mapGameToResp(g))
	}

	web.Respond(c, resps, http.StatusOK)
}

// CreateGame godoc
// @Summary Create game
// @Description creates new game
// @ID create-game
// @Accept  json
// @Produce json
// @Param   game body CreateGameReq true "create game"
// @Success 201 {object} IDResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [post]
func (g *Game) CreateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.creategame")
	defer span.End()

	var cg CreateGameReq
	err := web.Decode(c, &cg)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding new game"))
		return
	}
	span.SetAttributes(attribute.String("data.name", cg.Name))

	claims, err := web.GetClaims(c)
	if err != nil {
		c.Error(errors.Wrap(err, "getting claims from context"))
		return
	}

	create := mapToCreateGame(&cg)
	create.Publisher = claims.Name
	id, err := g.Storage.CreateGame(ctx, create)
	if err != nil {
		c.Error(errors.Wrap(err, "adding new game"))
		return
	}

	web.Respond(c, IDResp{ID: id}, http.StatusCreated)
}

// UpdateGame godoc
// @Summary Update game
// @Description updates game by ID
// @ID update-game
// @Accept  json
// @Produce json
// @Param  	id   path int32 		true "Game ID"
// @Param  	game body UpdateGameReq true "update game"
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
		c.Error(err)
		return
	}
	var ugr UpdateGameReq
	if err := web.Decode(c, &ugr); err != nil {
		c.Error(errors.Wrap(err, "decoding game update"))
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(id)))

	game, err := g.Storage.GetGameByID(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "retrieve game with id %d", id))
		return
	}

	update := mapToUpdateGame(game, ugr)
	err = g.Storage.UpdateGame(ctx, id, update)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "updating game with id %v", id))
		return
	}

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
		c.Error(err)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(id)))

	err = g.Storage.DeleteGame(ctx, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "deleting game with id %v", id))
		return
	}

	web.Respond(c, nil, http.StatusNoContent)
}
