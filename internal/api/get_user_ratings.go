package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetUserRatings returns user ratings for specified games
func (p *Provider) GetUserRatings(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getUserRatings")
	defer span.End()

	var ur GetUserRatingsRequest
	if err := p.decodeRequest(r, &ur); err != nil {
		web.RespondError(w, err)
		return
	}
	idsVal := make([]int, 0, len(ur.GameIDs))
	for _, v := range ur.GameIDs {
		idsVal = append(idsVal, int(v))
	}
	span.SetAttributes(attribute.IntSlice("data.ids", idsVal))

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	userID := claims.UserID()

	ratings, err := p.gameFacade.GetUserRatings(ctx, userID)
	if err != nil {
		p.respondError(w, "get user ratings", err, zap.String("user_id", userID))
		return
	}

	userRatings := make(UserRatingsResponse)
	for _, id := range ur.GameIDs {
		if rating, ok := ratings[id]; ok {
			userRatings[id] = rating
		}
	}

	web.Respond(w, userRatings, http.StatusOK)
}
