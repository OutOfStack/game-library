package api

import (
	"errors"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
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
func (p *Provider) GetGenres(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.getGenres")
	defer span.End()

	list, err := p.gameFacade.GetGenres(ctx)
	if err != nil {
		p.log.Error("get genres", zap.Error(err))
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
