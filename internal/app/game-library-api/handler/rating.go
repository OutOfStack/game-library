package handler

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
)

// RateGame godoc
// @Summary Rate game
// @Description rates game
// @ID rate-game
// @Accept  json
// @Produce json
// @Param   id 		path int64 				true "game ID"
// @Param	rating 	body CreateRatingReq 	true "game rating"
// @Success 200 {object} RatingResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/rate [post]
func (g *Game) RateGame(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.rating.rategame")
	defer span.End()

	gameID, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	var cr CreateRatingReq
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

	_, err = repo.Retrieve(ctx, g.DB, gameID)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "retrieve game"))
		return
	}

	rating := mapToCreateRating(&cr, gameID, userID)
	err = repo.AddRating(ctx, g.DB, rating)
	if err != nil {
		c.Error(errors.Wrap(err, "rate game"))
		return
	}

	resp := mapToRatingResp(rating)
	web.Respond(c, resp, http.StatusOK)
}

// GetUserRatings godoc
// @Summary Get user ratings for specified games
// @Description returns user ratings for specified games
// @ID get-user-ratings
// @Produce json
// @Param   gameIds body UserRatingsReq 	true "games ids"
// @Success 200 {object} map[int64]uint8
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /user/ratings [post]
func (g *Game) GetUserRatings(c *gin.Context) {
	ctx, span := trace.SpanFromContext(c.Request.Context()).Tracer().Start(c.Request.Context(), "handlers.rating.getuserratings")
	defer span.End()

	var ur UserRatingsReq
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

	ratings, err := repo.GetUserRatings(ctx, g.DB, userID, ur.GameIDs)

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
