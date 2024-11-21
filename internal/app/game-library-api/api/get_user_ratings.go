package api

import (
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetUserRatings godoc
// @Summary Get user ratings for specified games
// @Description returns user ratings for specified games
// @ID get-user-ratings
// @Produce json
// @Param   gameIds body api.GetUserRatingsRequest 	true "games ids"
// @Success 200 {object} map[int32]uint8
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /user/ratings [post]
func (p *Provider) GetUserRatings(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "api.getUserRatings")
	defer span.End()

	var ur api.GetUserRatingsRequest
	err := web.Decode(r, &ur)
	if err != nil {
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
		p.log.Error("get claims from ctx", zap.Error(err))
		web.Respond500(w)
		return
	}

	userID := claims.UserID()

	ratings, err := p.gameFacade.GetUserRatings(ctx, userID)
	if err != nil {
		p.log.Error("get user ratings", zap.String("user_id", userID), zap.Error(err))
		web.Respond500(w)
		return
	}

	userRatings := make(map[int32]uint8)
	for _, id := range ur.GameIDs {
		if r, ok := ratings[id]; ok {
			userRatings[id] = r
		}
	}

	web.Respond(w, userRatings, http.StatusOK)
}
