package facade

import (
	"fmt"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"golang.org/x/net/context"
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
	genres, err := p.GetGenres(ctx)
	if err != nil {
		return model.Genre{}, fmt.Errorf("get genres: %v", err)
	}

	for _, genre := range genres {
		if genre.ID == id {
			return genre, nil
		}
	}

	return model.Genre{}, apperr.NewNotFoundError("genre", id)
}
