package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/api/model"
	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/web"
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
func (p *Provider) GetTopCompanies(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getTopCompanies")
	defer span.End()

	companyType := r.URL.Query().Get("type")
	if companyType != model.CompanyTypeDeveloper && companyType != model.CompanyTypePublisher {
		web.RespondError(w, web.NewErrorFromMessage("invalid company type: should be one of [dev, pub]", http.StatusBadRequest))
		return
	}

	span.SetAttributes(att.String("type", companyType))

	list, err := p.gameFacade.GetTopCompanies(ctx, companyType, topCompaniesLimit)
	if err != nil {
		p.log.Error("get top companies", zap.String("type", companyType), zap.Error(err))
		web.Respond500(w)
		return
	}

	resp := make([]api.Company, 0, len(list))
	for _, company := range list {
		resp = append(resp, api.Company{
			ID:   company.ID,
			Name: company.Name,
		})
	}

	web.Respond(w, resp, http.StatusOK)
}
