package repo

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

// AddRating adds rating to game
func (s *Storage) AddRating(ctx context.Context, cr CreateRating) error {
	ctx, span := tracer.Start(ctx, "db.rating.add")
	defer span.End()

	const q = `
	INSERT INTO ratings
	(game_id, user_id, rating)
	VALUES ($1, $2, $3)
	ON CONFLICT (game_id, user_id) DO UPDATE SET rating = $3`

	if _, err := s.db.ExecContext(ctx, q, cr.GameID, cr.UserID, cr.Rating); err != nil {
		return fmt.Errorf("adding ratings to game with id %v from user with id %v: %w", cr.GameID, cr.UserID, err)
	}

	return nil
}

// RemoveRating removes rating of game
func (s *Storage) RemoveRating(ctx context.Context, rr RemoveRating) error {
	ctx, span := tracer.Start(ctx, "db.rating.remove")
	defer span.End()

	const q = `
	DELETE FROM ratings
	WHERE game_id = $1 AND user_id = $2`

	if _, err := s.db.ExecContext(ctx, q, rr.GameID, rr.UserID); err != nil {
		return fmt.Errorf("remove rating of game with id %v from user with id %v: %w", rr.GameID, rr.UserID, err)
	}

	return nil
}

// GetUserRatingsByGamesIDs returns user ratings for specified games
func (s *Storage) GetUserRatingsByGamesIDs(ctx context.Context, userID string, gameIDs []int32) ([]UserRating, error) {
	ctx, span := tracer.Start(ctx, "db.rating.getuserratingsbyids")
	defer span.End()

	ratings := make([]UserRating, 0)
	const q = `
	SELECT game_id, rating, user_id
	FROM ratings
	WHERE user_id = $1 AND game_id = ANY($2)`

	if err := s.db.SelectContext(ctx, &ratings, q, userID, pq.Array(gameIDs)); err != nil {
		return nil, err
	}

	return ratings, nil
}

// GetUserRatings returns all user ratings
func (s *Storage) GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error) {
	ctx, span := tracer.Start(ctx, "db.rating.getuserratings")
	defer span.End()

	ratings := make([]UserRating, 0)
	const q = `
	SELECT game_id, rating, user_id
	FROM ratings
	WHERE user_id = $1`

	if err := s.db.SelectContext(ctx, &ratings, q, userID); err != nil {
		return nil, err
	}

	userRatings := make(map[int32]uint8, len(ratings))
	for _, r := range ratings {
		userRatings[r.GameID] = r.Rating
	}

	return userRatings, nil
}
