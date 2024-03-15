package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/gin-gonic/gin"
	att "go.opentelemetry.io/otel/attribute"
)

const (
	topCompaniesLimit = 10

	companyTypeDeveloper = "dev"
	companyTypePublisher = "pub"
)

// GetTopCompanies godoc
// @Summary Get top companies
// @Description returns top companies based on type (developer or publisher)
// @ID get-top-companies
// @Produce json
// @Param   type query string true "company type (dev or pub)" Enums(pub, dev)
// @Success 200 {array}  Company
// @Failure 400 {object} web.ErrorResponse "Invalid or missing company type"
// @Failure 500 {object} web.ErrorResponse
// @Router /companies/top [get]
func (p *Provider) GetTopCompanies(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getTopCompanies")
	defer span.End()

	companyType := c.Query("type")
	if companyType != companyTypeDeveloper && companyType != companyTypePublisher {
		web.Err(c, web.NewRequestError(errors.New("invalid company type: should be one of [dev, pub]"), http.StatusBadRequest))
		return
	}

	span.SetAttributes(att.String("type", companyType))

	list := make([]repo.Company, 0)
	err := cache.Get(ctx, p.cache, getTopCompaniesKey(companyType, topCompaniesLimit), &list, func() ([]repo.Company, error) {
		switch companyType {
		case companyTypeDeveloper:
			return p.storage.GetTopDevelopers(ctx, topCompaniesLimit)
		case companyTypePublisher:
			return p.storage.GetTopPublishers(ctx, topCompaniesLimit)
		}
		return nil, fmt.Errorf("unsopported companyType: %s", companyType)
	}, 0)
	if err != nil {
		web.Err(c, fmt.Errorf("get top companies: %v", err))
		return
	}

	resp := make([]Company, 0, len(list))
	for _, company := range list {
		resp = append(resp, Company{
			ID:   company.ID,
			Name: company.Name,
		})
	}

	web.Respond(c, resp, http.StatusOK)
}
