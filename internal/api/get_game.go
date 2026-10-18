package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetGame returns game by ID
func (p *Provider) GetGame(w http.ResponseWriter, r *http.Request, id GameID) {
	ctx, span := tracer.Start(r.Context(), "getGame")
	defer span.End()

	if id <= 0 {
		web.RespondError(w, errInvalidID)
		return
	}

	span.SetAttributes(att.Int("data.id", int(id)))

	details, err := p.gameFacade.GetGameDetailsByID(ctx, id)
	if err != nil {
		p.respondError(w, "get game", err, zap.Int32("id", id))
		return
	}

	resp, err := p.mapToGameDetailsResponse(ctx, details)
	if err != nil {
		p.respondError(w, "map game to response", err, zap.Int32("id", id))
		return
	}

	web.Respond(w, resp, http.StatusOK)
}
