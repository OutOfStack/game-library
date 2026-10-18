package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/web"
)

const (
	topGenresLimit = 10
)

// GetTopGenres returns top genres based on amount of games having it
func (p *Provider) GetTopGenres(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getTopGenres")
	defer span.End()

	list, err := p.gameFacade.GetTopGenres(ctx, topGenresLimit)
	if err != nil {
		p.respondError(w, "get top genres", err)
		return
	}

	resp := make([]Genre, 0, len(list))
	for _, genre := range list {
		resp = append(resp, Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	web.Respond(w, resp, http.StatusOK)
}
