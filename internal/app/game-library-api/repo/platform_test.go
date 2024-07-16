package repo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

func TestGetPlatforms_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := context.Background()

	platform := model.Platform{
		Name:         td.String(),
		Abbreviation: td.String(),
		IGDBID:       td.Int64(),
	}

	id, err := createPlatform(ctx, platform)
	require.NoError(t, err)

	platforms, err := s.GetPlatforms(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(platforms), 1, "platforms len should be greater than 1")

	want := platform
	var got model.Platform
	for _, p := range platforms {
		if p.ID == id {
			got = p
			break
		}
	}

	require.Equal(t, id, got.ID, "id should be equal")
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.Abbreviation, got.Abbreviation, "abbreviation should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

func TestPlatformByID_PlatformExists_ShouldReturnPlatform(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := context.Background()

	platform := model.Platform{
		Name:         td.String(),
		Abbreviation: td.String(),
		IGDBID:       td.Int64(),
	}

	id, err := createPlatform(ctx, platform)
	require.NoError(t, err)

	gotPlatform, err := s.GetPlatformByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, gotPlatform.ID, "id should be equal")
	require.Equal(t, platform.Name, gotPlatform.Name, "name should be equal")
	require.Equal(t, platform.Abbreviation, gotPlatform.Abbreviation, "abbreviation should be equal")
	require.Equal(t, platform.IGDBID, gotPlatform.IGDBID, "igdb id should be equal")
}

func TestGetPlatformByID_PlatformNotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := context.Background()

	platform := model.Platform{
		Name:         td.String(),
		Abbreviation: td.String(),
		IGDBID:       td.Int64(),
	}

	_, err := createPlatform(ctx, platform)
	require.NoError(t, err)

	randomID := td.Int32()
	gotPlatform, err := s.GetPlatformByID(ctx, randomID)
	require.ErrorIs(t, err, apperr.NewNotFoundError("platform", randomID), "err should be NotFound")
	require.Zero(t, gotPlatform.ID, "got id should be 0")
}

func createPlatform(ctx context.Context, p model.Platform) (id int32, err error) {
	const q = `
		INSERT INTO platforms (name, abbreviation, igdb_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (igdb_id) DO NOTHING
		RETURNING id`

	if err = db.QueryRowContext(ctx, q, p.Name, p.Abbreviation, p.IGDBID).Scan(&id); err != nil {
		return 0, fmt.Errorf("create platform: %v", err)
	}

	return id, nil
}
