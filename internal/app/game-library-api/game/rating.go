package game

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// AddRating adds rating to game
func AddRating(ctx context.Context, db *sqlx.DB, cr CreateRating, userID string) (*Rating, error) {
	_, err := Retrieve(ctx, db, cr.GameID)
	if err != nil {
		return nil, err
	}

	const q = `
	insert into ratings
	(game_id, user_id, rating)
	values ($1, $2, $3)
	on conflict (game_id, user_id) do update set rating = $3`

	_, err = db.ExecContext(ctx, q, cr.GameID, userID, cr.Rating)
	if err != nil {
		return nil, fmt.Errorf("adding ratings to game with id %v from user with id %v: %w", cr.GameID, userID, err)
	}

	rating := &Rating{
		GameID: cr.GameID,
		UserID: userID,
		Rating: cr.Rating,
	}

	return rating, nil
}

// GetUserRatings returns ratings of user for specified games
func GetUserRatings(ctx context.Context, db *sqlx.DB, userId string, gameIDs []int64) ([]UserRating, error) {
	ratings := []UserRating{}
	const q = `
	select game_id, rating 
	from ratings
	where user_id = $1 and game_id = ANY($2)`

	if err := db.SelectContext(ctx, &ratings, q, userId, pq.Array(gameIDs)); err != nil {
		return nil, err
	}

	return ratings, nil
}
