package facade

import (
	"context"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"go.uber.org/zap"
)

// RateGame sets rating to game by user
func (p *Provider) RateGame(ctx context.Context, gameID int32, userID string, rating uint8) error {
	_, err := p.storage.GetGameByID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game %d by id: %w", gameID, err)
	}

	if rating == 0 {
		err = p.storage.RemoveRating(ctx, model.RemoveRating{
			UserID: userID,
			GameID: gameID,
		})
	} else {
		err = p.storage.AddRating(ctx, model.CreateRating{
			Rating: rating,
			UserID: userID,
			GameID: gameID,
		})
	}
	if err != nil {
		return fmt.Errorf("set rating %d to game %d by user %s: %w", rating, gameID, userID, err)
	}

	// update avg game rating
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()

		gErr := p.storage.UpdateGameRating(bCtx, gameID)
		if gErr != nil {
			p.log.Error("update game rating", zap.Int32("id", gameID), zap.Error(gErr))
		}

		// invalidate game cache
		gErr = cache.Delete(bCtx, p.cache, getGameKey(gameID))
		if gErr != nil {
			p.log.Error("remove game cache", zap.Int32("id", gameID), zap.Error(gErr))
		}
		// recache game
		_, gErr = p.GetGameByID(bCtx, gameID)
		if gErr != nil {
			p.log.Error("recache game", zap.Int32("id", gameID), zap.Error(gErr))
		}
	}()

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()

		// invalidate user ratings
		key := getUserRatingsKey(userID)
		gErr := cache.Delete(bCtx, p.cache, key)
		if gErr != nil {
			p.log.Error("remove cache by key", zap.String("key", key), zap.Error(gErr))
		}
		// recache user ratings
		_, gErr = p.GetUserRatings(bCtx, userID)
		if gErr != nil {
			p.log.Error("recache user ratings", zap.String("user_id", userID), zap.Error(gErr))
		}
	}()

	return nil
}

// GetUserRatings returns user's rating for specified games
func (p *Provider) GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error) {
	list := make(map[int32]uint8)
	err := cache.Get(ctx, p.cache, getUserRatingsKey(userID), &list, func() (map[int32]uint8, error) {
		return p.storage.GetUserRatings(ctx, userID)
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("get user %s games ratings: %v", userID, err)
	}

	return list, nil
}
