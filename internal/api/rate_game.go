package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// RateGame rates game
func (p *Provider) RateGame(w http.ResponseWriter, r *http.Request, id GameID) {
	ctx, span := tracer.Start(r.Context(), "rateGame")
	defer span.End()

	if id <= 0 {
		web.RespondError(w, errInvalidID)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(id)))

	var cr CreateRatingRequest
	if err := p.decodeRequest(r, &cr); err != nil {
		web.RespondError(w, err)
		return
	}

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	userID := claims.UserID()

	err = p.gameFacade.RateGame(ctx, id, userID, cr.Rating)
	if err != nil {
		p.respondError(w, "rate game", err, zap.Int32("id", id), zap.String("user_id", userID))
		return
	}

	web.Respond(w, RatingResponse{
		GameID: id,
		Rating: cr.Rating,
	}, http.StatusOK)
}
