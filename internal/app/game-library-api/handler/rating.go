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

// GetUserRatings godoc
// @Summary Get user ratings for specified games
// @Description returns user ratings for specified games
// @ID get-user-ratings
// @Produce json
// @Param   gameIds body game.UserRatings true "games ids"
// @Success 200 {object} map[int64]uint8
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /user/ratings [post]
func (g *Game) GetUserRatings(c *gin.Context) {
	var ur repo.UserRatings
	err := web.Decode(c, &ur)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding user ratings"))
		return
	}

	claims, err := web.GetClaims(c)
	if err != nil {
		c.Error(errors.Wrap(err, "getting claims from context"))
		return
	}

	userID := claims.Subject

	ratings, err := repo.GetUserRatings(c, g.DB, userID, ur.GameIDs)

	if err != nil {
		c.Error(errors.Wrap(err, "getting user ratings"))
		return
	}

	userRatings := make(map[int64]uint8)
	for _, r := range ratings {
		userRatings[r.GameID] = r.Rating
	}

	web.Respond(c, userRatings, http.StatusOK)
}
