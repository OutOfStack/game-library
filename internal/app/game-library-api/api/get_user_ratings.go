package api

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
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
func (p *Provider) GetUserRatings(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "api.getUserRatings")
	defer span.End()

	var ur api.GetUserRatingsRequest
	err := web.Decode(c, &ur)
	if err != nil {
		web.Err(c, fmt.Errorf("decoding user ratings: %w", err))
		return
	}
	idsVal := make([]int, 0, len(ur.GameIDs))
	for _, v := range ur.GameIDs {
		idsVal = append(idsVal, int(v))
	}
	span.SetAttributes(attribute.IntSlice("data.ids", idsVal))

	claims, err := web.GetClaims(c)
	if err != nil {
		web.Err(c, fmt.Errorf("getting claims from context: %w", err))
		return
	}

	userID := claims.UserID()

	ratings, err := p.gameFacade.GetUserRatings(ctx, userID)
	if err != nil {
		p.log.Error("get user ratings", zap.String("user_id", userID), zap.Error(err))
		web.Err(c, errors.New("internal error"))
		return
	}

	userRatings := make(map[int32]uint8)
	for _, id := range ur.GameIDs {
		if r, ok := ratings[id]; ok {
			userRatings[id] = r
		}
	}

	web.Respond(c, userRatings, http.StatusOK)
}
