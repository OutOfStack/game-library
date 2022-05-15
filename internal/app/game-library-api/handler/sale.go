package handler

import (
	"net/http"

	repo "github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
)

// AddSale godoc
// @Summary Create sale
// @Description Creates new sale
// @ID create-sale
// @Accept  json
// @Produce json
// @Param  	sale body 	 game.CreateSaleReq true "create sale"
// @Success 201 {object} game.SaleResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [post]
func (g *Game) AddSale(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.sale.addsale")
	defer span.End()

	var cs repo.CreateSaleReq
	err := web.Decode(c, &cs)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding new sale"))
		return
	}

	saleID, err := repo.AddSale(ctx, g.DB, cs)
	if err != nil {
		c.Error(errors.Wrapf(err, "adding new sale"))
		return
	}

	getSale := cs.MapToSaleResp(saleID)

	web.Respond(c, getSale, http.StatusCreated)
}

// ListSales godoc
// @Summary List all sales
// @Description Returns all sales
// @ID get-sales
// @Produce json
// @Success 200 {array}  game.SaleResp
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [get]
func (g *Game) ListSales(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.sale.listsales")
	defer span.End()

	list, err := repo.GetSales(ctx, g.DB)
	if err != nil {
		c.Error(errors.Wrap(err, "getting sales list"))
		return
	}

	getSales := []repo.SaleResp{}
	for _, s := range list {
		getSales = append(getSales, *s.MapToSaleResp())
	}

	web.Respond(c, getSales, http.StatusOK)
}
