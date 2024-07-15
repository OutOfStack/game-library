package api

import (
	"errors"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
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
func (p *Provider) GetPlatforms(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.getPlatforms")
	defer span.End()

	list, err := p.gameFacade.GetPlatforms(ctx)
	if err != nil {
		p.log.Error("get platforms", zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	resp := make([]api.Platform, 0, len(list))
	for _, p := range list {
		resp = append(resp, api.Platform{
			ID:           p.ID,
			Name:         p.Name,
			Abbreviation: p.Abbreviation,
		})
	}

	web.Respond(c, resp, http.StatusOK)
}
