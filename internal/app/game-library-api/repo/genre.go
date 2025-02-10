package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/georgysavva/scany/v2/pgxscan"
)

// CreateGenre creates new genre
func (s *Storage) CreateGenre(ctx context.Context, g model.Genre) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.createGenre")
	defer span.End()

	const q = `
		INSERT INTO genres (name, igdb_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (igdb_id) DO NOTHING
		RETURNING id`

	if err = s.db.QueryRow(ctx, q, g.Name, g.IGDBID, time.Now()).Scan(&id); err != nil {
		return 0, fmt.Errorf("create genre with name %s and igdb id %d: %v", g.Name, g.IGDBID, err)
	}

	return id, nil
}

// GetGenres returns genres
func (s *Storage) GetGenres(ctx context.Context) (genres []model.Genre, err error) {
	ctx, span := tracer.Start(ctx, "db.getGenres")
	defer span.End()

	const q = `
		SELECT id, name, igdb_id
		FROM genres`

	if err = pgxscan.Select(ctx, s.db, &genres, q); err != nil {
		return nil, err
	}

	return genres, nil
}

// GetGenreByID returns genre by id
// If company does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetGenreByID(ctx context.Context, id int32) (genre model.Genre, err error) {
	ctx, span := tracer.Start(ctx, "db.getGenreByID")
	defer span.End()

	const q = `
		SELECT id, name, igdb_id
		FROM genres
		WHERE id = $1`

	if err = pgxscan.Get(ctx, s.db, &genre, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Genre{}, apperr.NewNotFoundError("genre", id)
		}
		return model.Genre{}, err
	}

	return genre, nil
}

// GetTopGenres returns genres
func (s *Storage) GetTopGenres(ctx context.Context, limit int64) (genres []model.Genre, err error) {
	ctx, span := tracer.Start(ctx, "db.getTopGenres")
	defer span.End()

	const q = `
		SELECT gr.id, gr.name, gr.igdb_id
		FROM genres gr
		JOIN (
			SELECT UNNEST(genres) AS genre_id FROM games
		) AS g ON gr.id = g.genre_id
		GROUP BY gr.id, gr.name, gr.igdb_id
		ORDER BY COUNT(*) DESC
		LIMIT $1`

	if err = pgxscan.Select(ctx, s.db, &genres, q, limit); err != nil {
		return nil, err
	}

	return genres, nil
}
