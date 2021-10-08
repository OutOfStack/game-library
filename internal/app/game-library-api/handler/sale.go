package handler

import (
	"net/http"

	repo "github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// AddSale godoc
// @Summary Create sale
// @Description Creates new sale
// @ID create-sale
// @Accept  json
// @Produce json
// @Param  	sale body 	 game.CreateSale true "create sale"
// @Success 201 {object} game.GetSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [post]
func (g *Game) AddSale(c *gin.Context) {
	var cs repo.CreateSale
	err := web.Decode(c, &cs)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding new sale"))
		return
	}

	saleId, err := repo.AddSale(c.Request.Context(), g.DB, cs)
	if err != nil {
		c.Error(errors.Wrapf(err, "adding new sale"))
		return
	}

	getSale := cs.MapToGetSale(saleId)

	web.Respond(c, getSale, http.StatusCreated)
}

// ListSales godoc
// @Summary List all sales
// @Description Returns all sales
// @ID get-sales
// @Produce json
// @Success 200 {array}  game.GetSale
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [get]
func (g *Game) ListSales(c *gin.Context) {
	list, err := repo.ListSales(c, g.DB)
	if err != nil {
		c.Error(errors.Wrap(err, "getting sales list"))
		return
	}

	getSales := []repo.GetSale{}
	for _, s := range list {
		getSales = append(getSales, *s.MapToGetSale())
	}

	web.Respond(c, getSales, http.StatusOK)
}
