package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/api/trace"
)

// AddRating adds rating to game
// If such entity does not exist returns error ErrNotFound{}
func AddRating(ctx context.Context, db *sqlx.DB, cr CreateRating) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.rating.addrating")
	defer span.End()

	const q = `
	insert into ratings
	(game_id, user_id, rating)
	values ($1, $2, $3)
	on conflict (game_id, user_id) do update set rating = $3`

	_, err := db.ExecContext(ctx, q, cr.GameID, cr.UserID, cr.Rating)
	if err != nil {
		return fmt.Errorf("adding ratings to game with id %v from user with id %v: %w", cr.GameID, cr.UserID, err)
	}

	return nil
}

// GetUserRatings returns ratings of user for specified games
func GetUserRatings(ctx context.Context, db *sqlx.DB, userID string, gameIDs []int64) ([]UserRating, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.rating.getuserratings")
	defer span.End()

	ratings := []UserRating{}
	const q = `
	select game_id, rating 
	from ratings
	where user_id = $1 and game_id = any($2)`

	if err := db.SelectContext(ctx, &ratings, q, userID, pq.Array(gameIDs)); err != nil {
		return nil, err
	}

	return ratings, nil
}
