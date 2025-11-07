package api

import (
	"net/http"
	"time"

	api "github.com/OutOfStack/game-library/internal/api/model"
	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/web"
	"go.uber.org/zap"
)

// GetGameModerations godoc
// @Summary Get game moderations
// @Description returns all moderation records for specified game id
// @Security BearerAuth
// @ID get-game-moderations
// @Produce json
// @Param   id path int32 true "Game ID"
// @Success 200 {array} api.ModerationItem
// @Failure 400 {object} web.ErrorResponse
// @Failure 403 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/moderations [get]
func (p *Provider) GetGameModerations(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getGameModerations")
	defer span.End()

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.log.Error("get claims from context", zap.Error(err))
		web.Respond500(w)
		return
	}

	id, err := web.GetIDParam(r)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	mods, err := p.gameFacade.GetGameModerations(ctx, id, claims.Name)
	if err != nil {
		if appErr, ok := apperr.IsAppError(err); ok {
			web.RespondError(w, web.NewError(appErr, appErr.HTTPStatusCode()))
			return
		}
		p.log.Error("get game moderations", zap.Int32("id", id), zap.Error(err))
		web.Respond500(w)
		return
	}

	resp := make([]api.ModerationItem, 0, len(mods))
	for _, m := range mods {
		it := api.ModerationItem{
			ID:      m.ID,
			Status:  m.Status,
			Details: m.Details,
		}
		if m.CreatedAt.Valid {
			it.CreatedAt = m.CreatedAt.Time.Format(time.RFC3339)
		}
		if m.UpdatedAt.Valid {
			it.UpdatedAt = m.UpdatedAt.Time.Format(time.RFC3339)
		}
		resp = append(resp, it)
	}

	web.Respond(w, resp, http.StatusOK)
}
