package repo

import (
	"context"
	"fmt"
)

// CreateGenre creates new genre
func (s *Storage) CreateGenre(ctx context.Context, g Genre) (int32, error) {
	ctx, span := tracer.Start(ctx, "db.genre.create")
	defer span.End()

	var id int32
	const q = `
	INSERT INTO genres
	(name, igdb_id)
	VALUES ($1, $2)
	ON CONFLICT (igdb_id) DO NOTHING
	RETURNING id`

	if err := s.db.QueryRowContext(ctx, q, g.Name, g.IGDBID).Scan(&id); err != nil {
		return 0, fmt.Errorf("create genre with name %s and igdb id %d: %v", g.Name, g.IGDBID, err)
	}

	return id, nil
}

// GetGenres returns genres
func (s *Storage) GetGenres(ctx context.Context) ([]Genre, error) {
	ctx, span := tracer.Start(ctx, "db.genre.get")
	defer span.End()

	genres := make([]Genre, 0)
	const q = `
	SELECT id, name, igdb_id
	FROM genres`

	if err := s.db.SelectContext(ctx, &genres, q); err != nil {
		return nil, err
	}

	return genres, nil
}
