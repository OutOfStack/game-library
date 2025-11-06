package facade

import (
	"context"
	"fmt"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
)

// GetPlatforms returns all platforms
func (p *Provider) GetPlatforms(ctx context.Context) ([]model.Platform, error) {
	list := make([]model.Platform, 0)
	err := cache.Get(ctx, p.cache, getPlatformsKey(), &list, func() ([]model.Platform, error) {
		return p.storage.GetPlatforms(ctx)
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("get platforms: %v", err)
	}

	return list, nil
}

// GetPlatformsMap returns all platforms map
func (p *Provider) GetPlatformsMap(ctx context.Context) (map[int32]model.Platform, error) {
	platforms, err := p.GetPlatforms(ctx)
	if err != nil {
		return nil, fmt.Errorf("get platforms: %v", err)
	}

	m := make(map[int32]model.Platform, len(platforms))
	for _, pl := range platforms {
		m[pl.ID] = pl
	}

	return m, nil
}

// GetPlatformByID returns platform by id
func (p *Provider) GetPlatformByID(ctx context.Context, id int32) (model.Platform, error) {
	platform, err := p.storage.GetPlatformByID(ctx, id)
	if err != nil {
		if apperr.IsStatusCode(err, apperr.NotFound) {
			return model.Platform{}, err
		}
		return model.Platform{}, fmt.Errorf("get platform by id %d: %v", id, err)
	}

	return platform, nil
}
