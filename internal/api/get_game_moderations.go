package api

import (
	"net/http"
	"time"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	"go.uber.org/zap"
)

// GetGameModerations returns all moderation records for specified game id
func (p *Provider) GetGameModerations(w http.ResponseWriter, r *http.Request, id GameID) {
	ctx, span := tracer.Start(r.Context(), "getGameModerations")
	defer span.End()

	if id <= 0 {
		web.RespondError(w, errInvalidID)
		return
	}

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	mods, err := p.gameFacade.GetGameModerations(ctx, id, claims.Name)
	if err != nil {
		p.respondError(w, "get game moderations", err, zap.Int32("id", id))
		return
	}

	resp := make([]ModerationItem, 0, len(mods))
	for _, m := range mods {
		it := ModerationItem{
			ID:           m.ID,
			ResultStatus: m.Status,
			Details:      m.Details,
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
