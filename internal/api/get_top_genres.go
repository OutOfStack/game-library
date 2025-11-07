package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/api/model"
	"github.com/OutOfStack/game-library/internal/web"
	"go.uber.org/zap"
)

const (
	topGenresLimit = 10
)

// GetTopGenres godoc
// @Summary Get top genres
// @Description returns top genres based on amount of games having it
// @ID get-top-genres
// @Produce json
// @Success 200 {array}  api.Genre
// @Failure 500 {object} web.ErrorResponse
// @Router /genres/top [get]
func (p *Provider) GetTopGenres(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getTopGenres")
	defer span.End()

	list, err := p.gameFacade.GetTopGenres(ctx, topGenresLimit)
	if err != nil {
		p.log.Error("get top genres", zap.Error(err))
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
