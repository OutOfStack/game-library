package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/api/model"
	"github.com/OutOfStack/game-library/internal/web"
	"go.uber.org/zap"
)

// GetGenres godoc
// @Summary Get genres
// @Description returns all genres
// @ID get-genres
// @Produce json
// @Success 200 {array}  api.Genre
// @Failure 500 {object} web.ErrorResponse
// @Router /genres [get]
func (p *Provider) GetGenres(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getGenres")
	defer span.End()

	list, err := p.gameFacade.GetGenres(ctx)
	if err != nil {
		p.log.Error("get genres", zap.Error(err))
		web.Respond500(w)
		return
	}

	resp := make([]api.Genre, 0, len(list))
	for _, genre := range list {
		resp = append(resp, api.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	web.Respond(w, resp, http.StatusOK)
}
