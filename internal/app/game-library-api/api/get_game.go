package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetGame godoc
// @Summary Get game
// @Description returns game by ID
// @ID get-game-by-id
// @Produce json
// @Param 	id  path int32 true "Game ID"
// @Success 200 {object} api.GameResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [get]
func (p *Provider) GetGame(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getGame")
	defer span.End()

	id, err := web.GetIDParam(r)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	game, err := p.gameFacade.GetGameByID(ctx, id)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("get game", zap.Int32("id", id), zap.Error(err))
		web.Respond500(w)
		return
	}

	var resp api.GameResponse
	resp, err = p.mapToGameResponse(ctx, game)
	if err != nil {
		web.RespondError(w, web.NewErrorFromMessage("error converting response", http.StatusInternalServerError))
		return
	}
	web.Respond(w, resp, http.StatusOK)
}
