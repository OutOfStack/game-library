package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"go.uber.org/zap"
)

// GetUserGames godoc
// @Summary Get user games
// @Description returns all games for current user-publisher
// @Security BearerAuth
// @ID get-user-games
// @Produce json
// @Success 200 {array} api.GameResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /user/games [get]
func (p *Provider) GetUserGames(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "getPublisherGames")
	defer span.End()

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.log.Error("get claims from context", zap.Error(err))
		web.Respond500(w)
		return
	}

	publisher := claims.Name

	games, err := p.gameFacade.GetPublisherGames(ctx, publisher)
	if err != nil {
		p.log.Error("get publisher games", zap.String("publisher", publisher), zap.Error(err))
		web.Respond500(w)
		return
	}

	resp := make([]api.GameResponse, 0, len(games))
	for _, g := range games {
		mapped, mErr := p.mapToGameResponse(ctx, g)
		if mErr != nil {
			p.log.Error("map to game response", zap.Int32("game_id", g.ID), zap.Error(mErr))
			web.Respond500(w)
			return
		}
		resp = append(resp, mapped)
	}

	web.Respond(w, resp, http.StatusOK)
}
