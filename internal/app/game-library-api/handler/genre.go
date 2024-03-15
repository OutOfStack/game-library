package handler

import (
	"fmt"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/gin-gonic/gin"
)

// GetGenres godoc
// @Summary Get genres
// @Description returns all genres
// @ID get-genres
// @Produce json
// @Success 200 {array}  Genre
// @Failure 500 {object} web.ErrorResponse
// @Router /genres [get]
func (p *Provider) GetGenres(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getGenres")
	defer span.End()

	list, err := p.storage.GetGenres(ctx)
	if err != nil {
		web.Err(c, fmt.Errorf("get genres: %v", err))
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

// GetTopGenres godoc
// @Summary Get top genres
// @Description returns top genres based on amount of games having it
// @ID get-top-genres
// @Produce json
// @Success 200 {array}  Genre
// @Failure 500 {object} web.ErrorResponse
// @Router /genres/top [get]
func (p *Provider) GetTopGenres(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getTopGenres")
	defer span.End()

	list := make([]repo.Genre, 0)
	err := cache.Get(ctx, p.cache, getTopGenresKey(topGenresLimit), &list, func() ([]repo.Genre, error) {
		return p.storage.GetTopGenres(ctx, topGenresLimit)
	}, 0)
	if err != nil {
		web.Err(c, fmt.Errorf("get top genres: %v", err))
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
