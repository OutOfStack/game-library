package handler

import (
	"context"
	"net/http"

	repo "github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// AddSale godoc
// @Summary Create a sale
// @Description Creates new sale
// @ID create-sale
// @Accept  json
// @Produce json
// @Param  	sale body 	 game.CreateSale true "create sale"
// @Success 201 {object} game.GetSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [post]
func (g *Game) AddSale(ctx context.Context, c *gin.Context) error {
	var cs repo.CreateSale
	err := web.Decode(c, &cs)
	if err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	saleId, err := repo.AddSale(c.Request.Context(), g.DB, cs)
	if err != nil {
		return errors.Wrapf(err, "adding new sale")
	}

	getSale := cs.MapToGetSale(saleId)

	return web.Respond(ctx, c, getSale, http.StatusCreated)
}

// ListSales godoc
// @Summary List all sales
// @Description Returns all sales
// @ID get-sales
// @Produce json
// @Success 200 {array}  game.GetSale
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [get]
func (g *Game) ListSales(ctx context.Context, c *gin.Context) error {
	list, err := repo.ListSales(c.Request.Context(), g.DB)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	getSales := []repo.GetSale{}
	for _, s := range list {
		getSales = append(getSales, *s.MapToGetSale())
	}

	return web.Respond(ctx, c, getSales, http.StatusOK)
}

// ListGameSales godoc
// @Summary List game sales
// @Description returns sales for specified game
// @ID get-sales-by-game-id
// @Produce json
// @Param 	id  path int64 true "Game ID"
// @Success 200 {array}  game.GetGameSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /sales/game/{id} [get]
func (g *Game) ListGameSales(ctx context.Context, c *gin.Context) error {
	gameId, err := getIdParam(c)
	if err != nil {
		return err
	}

	gameSales, err := repo.ListGameSales(c.Request.Context(), g.DB, gameId)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "retrieving sales for game")
	}

	getGameSales := []repo.GetGameSale{}
	for _, gs := range gameSales {
		getGameSales = append(getGameSales, *gs.MapToGetGameSale())
	}

	return web.Respond(ctx, c, getGameSales, http.StatusOK)
}
