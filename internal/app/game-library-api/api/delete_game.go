package api

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/gin-gonic/gin"
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
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [delete]
func (p *Provider) DeleteGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.deleteGame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	err = p.gameFacade.DeleteGame(ctx, id)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.Err(c, web.NewRequestError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("delete game", zap.Int32("id", id), zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	web.Respond(c, nil, http.StatusNoContent)
}
