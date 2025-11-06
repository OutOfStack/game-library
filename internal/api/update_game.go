package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/api/model"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// UpdateGame godoc
// @Summary Update game
// @Description updates game by ID
// @Security BearerAuth
// @ID update-game
// @Accept  json
// @Produce json
// @Param  	id   path int32 				true "Game ID"
// @Param  	game body api.UpdateGameRequest true "update game"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 403 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [patch]
func (p *Provider) UpdateGame(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "updateGame")
	defer span.End()

	id, err := web.GetIDParam(r)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	var ugr api.UpdateGameRequest
	if err = p.decoder.Decode(r, &ugr); err != nil {
		web.RespondError(w, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.log.Error("get claims from context", zap.Error(err))
		web.Respond500(w)
		return
	}

	publisher := claims.Name

	update := mapToUpdateGame(&ugr, publisher)

	err = p.gameFacade.UpdateGame(ctx, id, update)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("update game", zap.Int32("id", id), zap.Error(err))
		web.Respond500(w)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}
