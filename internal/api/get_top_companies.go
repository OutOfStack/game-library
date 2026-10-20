package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	topCompaniesLimit = 10
)

// GetTopCompanies returns top companies based on amount of games having it
func (p *Provider) GetTopCompanies(w http.ResponseWriter, r *http.Request, params GetTopCompaniesParams) {
	ctx, span := tracer.Start(r.Context(), "getTopCompanies")
	defer span.End()

	companyType := params.Type
	if !companyType.Valid() {
		web.RespondError(w, errInvalidCompanyType)
		return
	}

	span.SetAttributes(att.String("type", string(companyType)))

	list, err := p.gameFacade.GetTopCompanies(ctx, string(companyType), topCompaniesLimit)
	if err != nil {
		p.respondError(w, "get top companies", err, zap.String("type", string(companyType)))
		return
	}

	resp := make([]Company, 0, len(list))
	for _, company := range list {
		resp = append(resp, Company{
			ID:   company.ID,
			Name: company.Name,
		})
	}

	web.Respond(w, resp, http.StatusOK)
}
