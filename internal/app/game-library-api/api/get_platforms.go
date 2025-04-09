package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"go.uber.org/zap"
)

// GetPlatforms godoc
// @Summary Get platforms
// @Description returns all platforms
// @ID get-platforms
// @Produce json
// @Success 200 {array}  api.Platform
// @Failure 500 {object} web.ErrorResponse
// @Router /platforms [get]
func (p *Provider) GetPlatforms(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getPlatforms")
	defer span.End()

	list, err := p.gameFacade.GetPlatforms(ctx)
	if err != nil {
		p.log.Error("get platforms", zap.Error(err))
		web.Respond500(w)
		return
	}

	resp := make([]api.Platform, 0, len(list))
	for _, pl := range list {
		resp = append(resp, api.Platform{
			ID:           pl.ID,
			Name:         pl.Name,
			Abbreviation: pl.Abbreviation,
		})
	}

	web.Respond(w, resp, http.StatusOK)
}
