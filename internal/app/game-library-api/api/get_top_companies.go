package api

import (
	"errors"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	topCompaniesLimit = 10
)

// GetTopCompanies godoc
// @Summary Get top companies
// @Description returns top companies based on amount of games having it
// @ID get-top-companies
// @Produce json
// @Param   type query string true "company type (dev or pub)" Enums(pub, dev)
// @Success 200 {array}  api.Company
// @Failure 400 {object} web.ErrorResponse "Invalid or missing company type"
// @Failure 500 {object} web.ErrorResponse
// @Router /companies/top [get]
func (p *Provider) GetTopCompanies(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.getTopCompanies")
	defer span.End()

	companyType := c.Query("type")
	if companyType != model.CompanyTypeDeveloper && companyType != model.CompanyTypePublisher {
		web.Err(c, web.NewRequestError(errors.New("invalid company type: should be one of [dev, pub]"), http.StatusBadRequest))
		return
	}

	span.SetAttributes(att.String("type", companyType))

	list, err := p.gameFacade.GetTopCompanies(ctx, companyType, topCompaniesLimit)
	if err != nil {
		p.log.Error("get top companies", zap.String("type", companyType), zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	resp := make([]api.Company, 0, len(list))
	for _, company := range list {
		resp = append(resp, api.Company{
			ID:   company.ID,
			Name: company.Name,
		})
	}

	web.Respond(c, resp, http.StatusOK)
}
