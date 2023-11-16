package handler

import (
	"fmt"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
)

// GetPlatforms godoc
// @Summary Get platforms
// @Description returns all platforms
// @ID get-platforms
// @Produce json
// @Success 200 {array}  Platform
// @Failure 500 {object} web.ErrorResponse
// @Router /platforms [get]
func (g *Game) GetPlatforms(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getPlatforms")
	defer span.End()

	list, err := g.storage.GetPlatforms(ctx)
	if err != nil {
		web.Err(c, fmt.Errorf("get platforms: %w", err))
		return
	}

	resp := make([]Platform, 0, len(list))
	for _, p := range list {
		resp = append(resp, Platform{
			ID:           p.ID,
			Name:         p.Name,
			Abbreviation: p.Abbreviation,
		})
	}

	web.Respond(c, resp, http.StatusOK)
}
