package facade

import (
	"context"
	"fmt"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
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

// GetPlatformByID returns platform by id
func (p *Provider) GetPlatformByID(ctx context.Context, id int32) (model.Platform, error) {
	platforms, err := p.GetPlatforms(ctx)
	if err != nil {
		return model.Platform{}, fmt.Errorf("get platforms: %v", err)
	}

	for _, platform := range platforms {
		if platform.ID == id {
			return platform, nil
		}
	}

	return model.Platform{}, apperr.NewNotFoundError("platform", id)
}
