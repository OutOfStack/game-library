package handler

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// GetGenres godoc
// @Summary Get genres
// @Description returns all genres
// @ID get-genres
// @Produce json
// @Success 200 {array}  Genre
// @Failure 500 {object} web.ErrorResponse
// @Router /genres [get]
func (g *Game) GetGenres(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getgenres")
	defer span.End()

	list, err := g.storage.GetGenres(ctx)
	if err != nil {
		c.Error(errors.Wrap(err, "get genres"))
		return
	}

	resp := make([]Genre, 0, len(list))
	for _, genre := range list {
		resp = append(resp, Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	web.Respond(c, resp, http.StatusOK)
}