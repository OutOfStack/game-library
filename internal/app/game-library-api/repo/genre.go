package repo

import (
	"context"
	"fmt"
	"time"
)

// CreateGenre creates new genre
func (s *Storage) CreateGenre(ctx context.Context, g Genre) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.createGenre")
	defer span.End()

	const q = `
	INSERT INTO genres
	(name, igdb_id, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (igdb_id) DO NOTHING
	RETURNING id`

	if err = s.db.QueryRowContext(ctx, q, g.Name, g.IGDBID, time.Now()).Scan(&id); err != nil {
		return 0, fmt.Errorf("create genre with name %s and igdb id %d: %v", g.Name, g.IGDBID, err)
	}

	return id, nil
}

// GetGenres returns genres
func (s *Storage) GetGenres(ctx context.Context) (genres []Genre, err error) {
	ctx, span := tracer.Start(ctx, "db.getGenres")
	defer span.End()

	const q = `
	SELECT id, name, igdb_id
	FROM genres`

	if err = s.db.SelectContext(ctx, &genres, q); err != nil {
		return nil, err
	}

	return genres, nil
}

// GetTopGenres returns genres
func (s *Storage) GetTopGenres(ctx context.Context, limit int64) (genres []Genre, err error) {
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

	if err = s.db.SelectContext(ctx, &genres, q, limit); err != nil {
		return nil, err
	}

	return genres, nil
}
