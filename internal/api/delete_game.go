package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// DeleteGame deletes game by ID
func (p *Provider) DeleteGame(w http.ResponseWriter, r *http.Request, id GameID) {
	ctx, span := tracer.Start(r.Context(), "deleteGame")
	defer span.End()

	if id <= 0 {
		web.RespondError(w, errInvalidID)
		return
	}

	span.SetAttributes(att.Int("data.id", int(id)))

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}
	publisher := claims.Name

	err = p.gameFacade.DeleteGame(ctx, id, publisher)
	if err != nil {
		p.respondError(w, "delete game", err, zap.Int32("id", id))
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}
