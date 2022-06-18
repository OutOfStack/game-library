package handler

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

// AddSale godoc
// @Summary Create sale
// @Description Creates new sale
// @ID create-sale
// @Accept  json
// @Produce json
// @Param  	sale body 	 CreateSaleReq true "create sale"
// @Success 201 {object} SaleResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [post]
func (g *Game) AddSale(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.sale.addsale")
	defer span.End()

	var cs CreateSaleReq
	err := web.Decode(c, &cs)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding new sale"))
		return
	}
	span.SetAttributes(attribute.String("data.name", cs.Name))

	sale := mapToCreateSale(&cs)
	saleID, err := g.Storage.AddSale(ctx, sale)
	if err != nil {
		c.Error(errors.Wrapf(err, "adding new sale"))
		return
	}

	resp := mapCreateSaleToSaleResp(sale, saleID)
	web.Respond(c, resp, http.StatusCreated)
}

// ListSales godoc
// @Summary List all sales
// @Description Returns all sales
// @ID get-sales
// @Produce json
// @Success 200 {array}  SaleResp
// @Failure 500 {object} web.ErrorResponse
// @Router /sales [get]
func (g *Game) ListSales(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.sale.listsales")
	defer span.End()

	list, err := g.Storage.GetSales(ctx)
	if err != nil {
		c.Error(errors.Wrap(err, "getting sales list"))
		return
	}

	resp := []SaleResp{}
	for _, s := range list {
		resp = append(resp, *mapSaleToSaleResp(&s))
	}

	web.Respond(c, resp, http.StatusOK)
}
