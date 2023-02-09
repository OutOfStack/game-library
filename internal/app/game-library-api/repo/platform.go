package repo

import (
	"context"
)

// GetPlatforms returns platforms
func (s *Storage) GetPlatforms(ctx context.Context) ([]Platform, error) {
	ctx, span := tracer.Start(ctx, "db.platform.get")
	defer span.End()

	platforms := make([]Platform, 0)
	const q = `
	SELECT id, name, abbreviation, igdb_id
	FROM platforms`

	if err := s.db.SelectContext(ctx, &platforms, q); err != nil {
		return nil, err
	}

	return platforms, nil
}
