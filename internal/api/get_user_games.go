package api

import (
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	"go.uber.org/zap"
)

// GetUserGames returns all games for current user-publisher
func (p *Provider) GetUserGames(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getPublisherGames")
	defer span.End()

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	publisher := claims.Name

	games, err := p.gameFacade.GetPublisherGames(ctx, publisher)
	if err != nil {
		p.respondError(w, "get publisher games", err, zap.String("publisher", publisher))
		return
	}

	resp := make([]GameResponse, 0, len(games))
	for _, g := range games {
		mapped, mErr := p.mapToGameResponse(ctx, g)
		if mErr != nil {
			p.respondError(w, "map to game response", mErr, zap.Int32("game_id", g.ID))
			return
		}
		resp = append(resp, mapped)
	}

	web.Respond(w, resp, http.StatusOK)
}
