package handler

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getplatforms")
	defer span.End()

	list, err := g.storage.GetPlatforms(ctx)
	if err != nil {
		c.Error(errors.Wrap(err, "get platforms"))
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
