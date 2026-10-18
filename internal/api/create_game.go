package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	att "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// CreateGame creates new game
func (p *Provider) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "createGame")
	defer span.End()

	var cg CreateGameRequest
	if err := p.decodeRequest(r, &cg); err != nil {
		web.RespondError(w, err)
		return
	}
	span.SetAttributes(att.String("data.name", cg.Name))

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	if claims.VerificationRequired {
		web.RespondError(w, web.NewErrorFromMessage("Email verification is required for adding games", http.StatusBadRequest))
		return
	}

	publisher := claims.Name

	create := mapToCreateGame(&cg, publisher)

	id, err := p.gameFacade.CreateGame(ctx, create)
	if err != nil {
		p.respondError(w, "create game", err, zap.String("name", create.Name), zap.String("user_id", claims.UserID()))
		return
	}

	web.Respond(w, IDResponse{ID: id}, http.StatusCreated)
}
