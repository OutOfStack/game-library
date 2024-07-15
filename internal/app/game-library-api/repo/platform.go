package repo

import (
	"context"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
)

// GetPlatforms returns platforms
func (s *Storage) GetPlatforms(ctx context.Context) (platforms []model.Platform, err error) {
	ctx, span := tracer.Start(ctx, "db.getPlatforms")
	defer span.End()

	const q = `
	SELECT id, name, abbreviation, igdb_id
	FROM platforms`

	if err = s.db.SelectContext(ctx, &platforms, q); err != nil {
		return nil, err
	}

	return platforms, nil
}
