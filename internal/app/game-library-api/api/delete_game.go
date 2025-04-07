package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// DeleteGame godoc
// @Summary Delete game
// @Description deletes game by ID
// @ID delete-game
// @Accept  json
// @Produce json
// @Param  	id path int32 true "Game ID"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 403 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [delete]
func (p *Provider) DeleteGame(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "deleteGame")
	defer span.End()

	id, err := web.GetIDParam(r)
	if err != nil {
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

	err = p.gameFacade.DeleteGame(ctx, id, publisher)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("delete game", zap.Int32("id", id), zap.Error(err))
		web.Respond500(w)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}
