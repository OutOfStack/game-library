package repo

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

// AddRating adds rating to game
// If such entity does not exist returns error ErrNotFound{}
func (s *Storage) AddRating(ctx context.Context, cr CreateRating) error {
	ctx, span := tracer.Start(ctx, "db.rating.addrating")
	defer span.End()

	const q = `
	insert into ratings
	(game_id, user_id, rating)
	values ($1, $2, $3)
	on conflict (game_id, user_id) do update set rating = $3`

	_, err := s.DB.ExecContext(ctx, q, cr.GameID, cr.UserID, cr.Rating)
	if err != nil {
		return fmt.Errorf("adding ratings to game with id %v from user with id %v: %w", cr.GameID, cr.UserID, err)
	}

	return nil
}

// GetUserRatings returns ratings of user for specified games
func (s *Storage) GetUserRatings(ctx context.Context, userID string, gameIDs []int64) ([]UserRating, error) {
	ctx, span := tracer.Start(ctx, "db.rating.getuserratings")
	defer span.End()

	ratings := []UserRating{}
	const q = `
	select game_id, rating 
	from ratings
	where user_id = $1 and game_id = any($2)`

	if err := s.DB.SelectContext(ctx, &ratings, q, userID, pq.Array(gameIDs)); err != nil {
		return nil, err
	}

	return ratings, nil
}
