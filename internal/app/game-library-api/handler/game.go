package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
)

// Game has handler methods for dealing with games
type Game struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// GetList godoc
// @Summary List games info
// @Description returns paginated games with extended properties
// @ID get-all-games-info
// @Produce json
// @Param pageSize query integer false "page size"
// @Param lastId   query integer false "last fetched Id"
// @Success 200 {array}  GameInfoResp
// @Failure 500 {object} web.ErrorResponse
// @Router /games [get]
func (g *Game) GetList(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.getlist")
	defer span.End()

	psParam := c.DefaultQuery("pageSize", "20")
	liParam := c.DefaultQuery("lastId", "0")
	pageSize, err := strconv.ParseInt(psParam, 10, 32)
	if err != nil || pageSize <= 0 {
		c.Error(web.NewRequestError(errors.New("Incorrect page size. Should be greater than 0"), http.StatusBadRequest))
		return
	}
	lastID, err := strconv.ParseInt(liParam, 10, 64)
	if err != nil || lastID < 0 {
		c.Error(web.NewRequestError(errors.New("Incorrect last Id. Should be greater or equal to 0"), http.StatusBadRequest))
		return
	}

	list, err := repo.GetInfos(ctx, g.DB, int(pageSize), lastID)

	if err != nil {
		c.Error(errors.Wrap(err, "getting games list"))
		return
	}

	resps := []GameInfoResp{}
	for _, g := range list {
		resps = append(resps, *mapToGameInfoResp(&g))
	}

	web.Respond(c, resps, http.StatusOK)
}

// Get godoc
// @Summary Get game info
// @Description returns game with extended properties by ID
// @ID get-game-info-by-id
// @Produce json
// @Param 	id  path int64 true "Game ID"
// @Success 200 {object} GameInfoResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [get]
func (g *Game) Get(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.get")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	game, err := repo.RetrieveInfo(ctx, g.DB, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "retrieving game with id %v", id))
		return
	}

	resp := mapToGameInfoResp(game)
	web.Respond(c, resp, http.StatusOK)
}

// Search godoc
// @Summary Searches games by part of name
// @Description returns list of games filtered by provided name
// @ID search-games-info
// @Produce json
// @Param name query string false "name to search by"
// @Success 200 {array}  GameInfoResp
// @Failure 500 {object} web.ErrorResponse
// @Router /games/search [get]
func (g *Game) Search(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.search")
	defer span.End()

	nameParam := c.DefaultQuery("name", "")
	if len(nameParam) < 2 {
		c.Error(web.NewRequestError(errors.New("Length of name to be searched by should be at least 2 characters"), http.StatusBadRequest))
		return
	}

	list, err := repo.SearchInfos(ctx, g.DB, nameParam)
	if err != nil {
		c.Error(errors.Wrap(err, "searching games list"))
		return
	}

	resps := []GameInfoResp{}
	for _, g := range list {
		resps = append(resps, *mapToGameInfoResp(&g))
	}

	web.Respond(c, resps, http.StatusOK)
}

// Create godoc
// @Summary Create game
// @Description creates new game
// @ID create-game
// @Accept  json
// @Produce json
// @Param   game body CreateGameReq true "create game"
// @Success 201 {object} GameResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [post]
func (g *Game) Create(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.create")
	defer span.End()

	var cg CreateGameReq
	err := web.Decode(c, &cg)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding new game"))
		return
	}

	claims, err := web.GetClaims(c)
	if err != nil {
		c.Error(errors.Wrap(err, "getting claims from context"))
		return
	}

	cg.Publisher = claims.Name

	create := mapToCreateGame(&cg)
	gameID, err := repo.Create(ctx, g.DB, create)
	if err != nil {
		c.Error(errors.Wrap(err, "adding new game"))
		return
	}

	resp := mapToGameResp(create, gameID)
	web.Respond(c, resp, http.StatusCreated)
}

// Update godoc
// @Summary Update game
// @Description updates game by ID
// @ID update-game-by-id
// @Accept  json
// @Produce json
// @Param  	id   path int64 		true "Game ID"
// @Param  	game body UpdateGameReq true "update game"
// @Success 200 {object} GameResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [patch]
func (g *Game) Update(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.update")
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

	game, err := repo.Retrieve(ctx, g.DB, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "retrieve game with id %v", id))
		return
	}

	update := mapToUpdateGame(game, &ugr)
	err = repo.Update(ctx, g.DB, update)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrapf(err, "updating game with id %v", id))
		return
	}

	resp := mapUpdateGameToGameResp(update)
	web.Respond(c, resp, http.StatusOK)
}

// Delete godoc
// @Summary Delete game
// @Description deletes game by ID
// @ID delete-game-by-id
// @Accept  json
// @Produce json
// @Param  	id path int64 true "Game ID"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [delete]
func (g *Game) Delete(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.delete")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = repo.Delete(ctx, g.DB, id)
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

// AddGameOnSale godoc
// @Summary Add game on sale
// @Description adds game on sale
// @ID add-game-on-sale
// @Accept  json
// @Produce json
// @Param  id 		path int64 				true "Game ID"
// @Param  gamesale body CreateGameSaleReq 	true "game sale"
// @Success 200 {object} GameSaleResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/sales [post]
func (g *Game) AddGameOnSale(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.addgameonsale")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	var cgs CreateGameSaleReq
	if err := web.Decode(c, &cgs); err != nil {
		c.Error(errors.Wrap(err, "decoding game sale"))
		return
	}

	_, err = repo.Retrieve(ctx, g.DB, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "get game"))
		return
	}

	sale, err := repo.RetrieveSale(ctx, g.DB, cgs.SaleID)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "add game on sale"))
		return
	}

	create := mapToCreateGameSale(&cgs, id)
	err = repo.AddGameOnSale(ctx, g.DB, create)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "add game on sale"))
		return
	}

	resp := mapToGameSaleResp(sale, create)
	web.Respond(c, resp, http.StatusOK)
}

// ListGameSales godoc
// @Summary List game sales
// @Description returns sales for specified game
// @ID get-game-sales-by-game-id
// @Produce json
// @Param 	id  path int64 true "Game ID"
// @Success 200 {array}  GameSaleResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/sales [get]
func (g *Game) ListGameSales(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.game.listgamesales")
	defer span.End()

	gameID, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	gameSales, err := repo.ListGameSales(ctx, g.DB, gameID)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "retrieving sales for game"))
		return
	}

	resps := []GameSaleResp{}
	for _, gs := range gameSales {
		resps = append(resps, *mapGameSaleToGameSaleResp(&gs))
	}

	web.Respond(c, resps, http.StatusOK)
}
