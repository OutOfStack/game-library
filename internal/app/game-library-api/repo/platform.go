package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/georgysavva/scany/v2/pgxscan"
)

// GetPlatforms returns platforms
func (s *Storage) GetPlatforms(ctx context.Context) (platforms []model.Platform, err error) {
	ctx, span := tracer.Start(ctx, "getPlatforms")
	defer span.End()

	const q = `
		SELECT id, name, abbreviation, igdb_id
		FROM platforms`

	if err = pgxscan.Select(ctx, s.db, &platforms, q); err != nil {
		return nil, err
	}

	return platforms, nil
}

// GetPlatformByID returns platform by id
// If company does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetPlatformByID(ctx context.Context, id int32) (platform model.Platform, err error) {
	ctx, span := tracer.Start(ctx, "getPlatformByID")
	defer span.End()

	const q = `
		SELECT id, name, abbreviation, igdb_id
		FROM platforms
		WHERE id = $1`

	if err = pgxscan.Get(ctx, s.db, &platform, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Platform{}, apperr.NewNotFoundError("platform", id)
		}
		return model.Platform{}, err
	}

	return platform, nil
}
