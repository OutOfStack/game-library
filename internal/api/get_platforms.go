package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/web"
)

// GetPlatforms returns all platforms
func (p *Provider) GetPlatforms(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getPlatforms")
	defer span.End()

	list, err := p.gameFacade.GetPlatforms(ctx)
	if err != nil {
		p.respondError(w, "get platforms", err)
		return
	}

	resp := make([]Platform, 0, len(list))
	for _, pl := range list {
		resp = append(resp, Platform{
			ID:           pl.ID,
			Name:         pl.Name,
			Abbreviation: pl.Abbreviation,
		})
	}

	web.Respond(w, resp, http.StatusOK)
}
