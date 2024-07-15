package api

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// RateGame godoc
// @Summary Rate game
// @Description rates game
// @ID rate-game
// @Accept  json
// @Produce json
// @Param   id 		path int32 					 true "game ID"
// @Param	rating 	body api.CreateRatingRequest true "game rating"
// @Success 200 {object} api.RatingResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/rate [post]
func (p *Provider) RateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.rateGame")
	defer span.End()

	gameID, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(gameID)))

	var cr api.CreateRatingRequest
	if err = web.Decode(c, &cr); err != nil {
		web.Err(c, fmt.Errorf("decoding rating: %w", err))
		return
	}

	claims, err := web.GetClaims(c)
	if err != nil {
		web.Err(c, fmt.Errorf("getting claims from context: %w", err))
		return
	}

	userID := claims.UserID()

	err = p.gameFacade.RateGame(ctx, gameID, userID, cr.Rating)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.Err(c, web.NewRequestError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("rate game", zap.Int32("id", gameID), zap.String("user_id", userID), zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	web.Respond(c, &api.RatingResponse{
		GameID: gameID,
		Rating: cr.Rating,
	}, http.StatusOK)
}
