package api

import (
	"errors"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
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
func (p *Provider) GetTopGenres(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.getTopGenres")
	defer span.End()

	list, err := p.gameFacade.GetTopGenres(ctx, topGenresLimit)
	if err != nil {
		p.log.Error("get top genres", zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	resp := make([]api.Genre, 0, len(list))
	for _, genre := range list {
		resp = append(resp, api.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	web.Respond(c, resp, http.StatusOK)
}
