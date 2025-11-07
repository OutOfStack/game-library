package facade

import (
	"context"
	"fmt"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
)

// GetGenres returns all genres
func (p *Provider) GetGenres(ctx context.Context) ([]model.Genre, error) {
	list := make([]model.Genre, 0)
	err := cache.Get(ctx, p.cache, getGenresKey(), &list, func() ([]model.Genre, error) {
		return p.storage.GetGenres(ctx)
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("get genres: %v", err)
	}

	return list, nil
}

// GetGenresMap returns all genres map
func (p *Provider) GetGenresMap(ctx context.Context) (map[int32]model.Genre, error) {
	genres, err := p.GetGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("get genres: %v", err)
	}

	m := make(map[int32]model.Genre, len(genres))
	for _, g := range genres {
		m[g.ID] = g
	}

	return m, nil
}

// GetTopGenres returns top genres
func (p *Provider) GetTopGenres(ctx context.Context, limit int64) ([]model.Genre, error) {
	list := make([]model.Genre, 0)
	err := cache.Get(ctx, p.cache, getTopGenresKey(limit), &list, func() ([]model.Genre, error) {
		return p.storage.GetTopGenres(ctx, limit)
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("get top genres: %v", err)
	}

	return list, nil
}

// GetGenreByID returns genre by id
func (p *Provider) GetGenreByID(ctx context.Context, id int32) (model.Genre, error) {
	genre, err := p.storage.GetGenreByID(ctx, id)
	if err != nil {
		if apperr.IsStatusCode(err, apperr.NotFound) {
			return model.Genre{}, err
		}
		return model.Genre{}, fmt.Errorf("get genre by id %d: %v", id, err)
	}

	return genre, nil
}
