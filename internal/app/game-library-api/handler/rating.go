package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// RateGame godoc
// @Summary Rate game
// @Description rates game
// @ID rate-game
// @Accept  json
// @Produce json
// @Param   id 		path int32 					true "game ID"
// @Param	rating 	body CreateRatingRequest 	true "game rating"
// @Success 200 {object} RatingResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/rate [post]
func (g *Game) RateGame(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.rateGame")
	defer span.End()

	gameID, err := web.GetIDParam(c)
	if err != nil {
		web.Err(c, err)
		return
	}
	span.SetAttributes(attribute.Int("data.id", int(gameID)))

	var cr CreateRatingRequest
	if err = web.Decode(c, &cr); err != nil {
		web.Err(c, fmt.Errorf("decoding rating: %w", err))
		return
	}

	claims, err := web.GetClaims(c)
	if err != nil {
		web.Err(c, fmt.Errorf("getting claims from context: %w", err))
		return
	}

	userID := claims.Subject

	_, err = g.storage.GetGameByID(ctx, gameID)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound[int32]{}) {
			web.Err(c, web.NewRequestError(err, http.StatusNotFound))
			return
		}
		web.Err(c, fmt.Errorf("get game by id: %w", err))
		return
	}

	if cr.Rating == 0 {
		err = g.storage.RemoveRating(ctx, repo.RemoveRating{
			UserID: userID,
			GameID: gameID,
		})
	} else {
		err = g.storage.AddRating(ctx, repo.CreateRating{
			Rating: cr.Rating,
			UserID: userID,
			GameID: gameID,
		})
	}
	if err != nil {
		web.Err(c, fmt.Errorf("rate game: %w", err))
		return
	}

	err = g.storage.UpdateGameRating(ctx, gameID)
	if err != nil {
		g.log.Error("updating game rating", zap.Int32("gameID", gameID), zap.Error(err))
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 1*time.Second)
		defer cancel()

		// invalidate user ratings
		key := getUserRatingsKey(userID)
		err = cache.Delete(bCtx, g.cache, key)
		if err != nil {
			g.log.Error("remove cache by key", zap.String("key", key), zap.Error(err))
		}
		// recache user ratings
		err = cache.Get(bCtx, g.cache, key, &map[int32]uint8{}, func() (map[int32]uint8, error) {
			return g.storage.GetUserRatings(bCtx, userID)
		}, 0)
		if err != nil {
			g.log.Error("recache user ratings", zap.String("user_id", userID), zap.Error(err))
		}
	}()

	web.Respond(c, &RatingResponse{
		GameID: gameID,
		Rating: cr.Rating,
	}, http.StatusOK)
}

// GetUserRatings godoc
// @Summary Get user ratings for specified games
// @Description returns user ratings for specified games
// @ID get-user-ratings
// @Produce json
// @Param   gameIds body GetUserRatingsRequest 	true "games ids"
// @Success 200 {object} map[int32]uint8
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /user/ratings [post]
func (g *Game) GetUserRatings(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handlers.getUserRatings")
	defer span.End()

	var ur GetUserRatingsRequest
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

	userID := claims.Subject

	var ratings map[int32]uint8
	err = cache.Get(ctx, g.cache, getUserRatingsKey(userID), &ratings, func() (map[int32]uint8, error) {
		return g.storage.GetUserRatings(ctx, userID)
	}, 0)
	if err != nil {
		web.Err(c, fmt.Errorf("getting user ratings: %w", err))
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
