package api

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/gin-gonic/gin"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// UpdateGame godoc
// @Summary Update game
// @Description updates game by ID
// @ID update-game
// @Accept  json
// @Produce json
// @Param  	id   path int32 				true "Game ID"
// @Param  	game body api.UpdateGameRequest true "update game"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [patch]
func (p *Provider) UpdateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.updateGame")
	defer span.End()

	id, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	var ugr api.UpdateGameRequest
	if err = web.Decode(c, &ugr); err != nil {
		web.Err(c, fmt.Errorf("decoding game update: %w", err))
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	err = p.gameFacade.UpdateGame(ctx, id, model.UpdatedGame(ugr))
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.Err(c, web.NewRequestError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("update game", zap.Int32("id", id), zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	web.Respond(c, nil, http.StatusNoContent)
}
