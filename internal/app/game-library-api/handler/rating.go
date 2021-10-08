package handler

import (
	"net/http"

	repo "github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// RateGame godoc
// @Summary Rate game
// @Description rates game
// @ID rate-game
// @Accept  json
// @Produce json
// @Param   rating body game.CreateRating true "game rating"
// @Success 200 {object} game.Rating
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/rate [post]
func (g *Game) RateGame(c *gin.Context) {
	var cr repo.CreateRating
	if err := web.Decode(c, &cr); err != nil {
		c.Error(errors.Wrap(err, "decoding rating"))
		return
	}

	claims, err := web.GetClaims(c)
	if err != nil {
		c.Error(errors.Wrap(err, "getting claims from context"))
		return
	}

	userID := claims.Subject
	rating, err := repo.AddRating(c, g.DB, cr, userID)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "rate game"))
		return
	}

	web.Respond(c, rating, http.StatusOK)
}
