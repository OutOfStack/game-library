package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/web"
)

// GetGenres returns all genres
func (p *Provider) GetGenres(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getGenres")
	defer span.End()

	list, err := p.gameFacade.GetGenres(ctx)
	if err != nil {
		p.respondError(w, "get genres", err)
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
