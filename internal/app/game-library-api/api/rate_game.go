package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// RateGame godoc
// @Summary Rate game
// @Description rates game
// @Security BearerAuth
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
func (p *Provider) RateGame(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "rateGame")
	defer span.End()

	gameID, err := web.GetIDParam(r)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(gameID)))

	var cr api.CreateRatingRequest
	if err = web.Decode(p.log, r, &cr); err != nil {
		web.RespondError(w, err)
		return
	}

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.log.Error("get cliams from ctx", zap.Error(err))
		web.Respond500(w)
		return
	}

	userID := claims.UserID()

	err = p.gameFacade.RateGame(ctx, gameID, userID, cr.Rating)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("rate game", zap.Int32("id", gameID), zap.String("user_id", userID), zap.Error(err))
		web.Respond500(w)
		return
	}

	web.Respond(w, &api.RatingResponse{
		GameID: gameID,
		Rating: cr.Rating,
	}, http.StatusOK)
}
