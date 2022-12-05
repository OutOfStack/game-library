package handler

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// RateGame godoc
// @Summary Rate game
// @Description rates game
// @ID rate-game
// @Accept  json
// @Produce json
// @Param   id 		path int32 				true "game ID"
// @Param	rating 	body CreateRatingReq 	true "game rating"
// @Success 200 {object} RatingResp
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/rate [post]
func (g *Game) RateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.rating.rategame")
	defer span.End()

	gameID, err := web.GetIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(gameID)))

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

	_, err = g.Storage.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			c.Error(web.NewRequestError(err, http.StatusNotFound))
			return
		}
		c.Error(errors.Wrap(err, "get game by id"))
		return
	}

	rating := mapToCreateRating(&cr, gameID, userID)
	err = g.Storage.AddRating(ctx, rating)
	if err != nil {
		c.Error(errors.Wrap(err, "rate game"))
		return
	}

	err = g.Storage.UpdateGameRating(ctx, gameID)
	if err != nil {
		g.Log.Error("updating game rating", zap.Int32("gameID", gameID), zap.Error(err))
	}

	web.Respond(c, mapCreateRatingToResp(rating), http.StatusOK)
}

// GetUserRatings godoc
// @Summary Get user ratings for specified games
// @Description returns user ratings for specified games
// @ID get-user-ratings
// @Produce json
// @Param   gameIds body UserRatingsReq 	true "games ids"
// @Success 200 {object} map[int32]uint8
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /user/ratings [post]
func (g *Game) GetUserRatings(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.rating.getuserratings")
	defer span.End()

	var ur UserRatingsReq
	err := web.Decode(c, &ur)
	if err != nil {
		c.Error(errors.Wrap(err, "decoding user ratings"))
		return
	}
	idsVal := make([]int, 0, len(ur.GameIDs))
	for _, v := range ur.GameIDs {
		idsVal = append(idsVal, int(v))
	}
	span.SetAttributes(attribute.IntSlice("data.ids", idsVal))

	claims, err := web.GetClaims(c)
	if err != nil {
		c.Error(errors.Wrap(err, "getting claims from context"))
		return
	}

	userID := claims.Subject

	ratings, err := g.Storage.GetUserRatings(ctx, userID, ur.GameIDs)
	if err != nil {
		c.Error(errors.Wrap(err, "getting user ratings"))
		return
	}

	userRatings := make(map[int32]uint8)
	for _, r := range ratings {
		userRatings[r.GameID] = r.Rating
	}

	web.Respond(c, userRatings, http.StatusOK)
}
