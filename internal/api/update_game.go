package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// UpdateGame updates game by ID
func (p *Provider) UpdateGame(w http.ResponseWriter, r *http.Request, id GameID) {
	ctx, span := tracer.Start(r.Context(), "updateGame")
	defer span.End()

	if id <= 0 {
		web.RespondError(w, errInvalidID)
		return
	}
	var ugr UpdateGameRequest
	if err := p.decodeRequest(r, &ugr); err != nil {
		web.RespondError(w, err)
		return
	}
	span.SetAttributes(att.Int("data.id", int(id)))

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	publisher := claims.Name

	update := mapToUpdateGame(&ugr, publisher)

	err = p.gameFacade.UpdateGame(ctx, id, update)
	if err != nil {
		p.respondError(w, "update game", err, zap.Int32("id", id))
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}
