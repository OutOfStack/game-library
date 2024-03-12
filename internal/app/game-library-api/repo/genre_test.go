package repo_test

import (
	"context"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

// TestCreateGenre_IGDBIDIsNull_ShouldBeNoError tests case when we add genre without igdb id, and there should be no error
func TestCreateGenre_IGDBIDIsNull_ShouldBeNoError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	genre := repo.Genre{
		Name: td.String(),
	}

	_, err := s.CreateGenre(context.Background(), genre)
	require.NoError(t, err)
}

// TestGetGenres_DataExists_ShouldBeEqual tests case when we add one genre, then fetch first genre, and they should be equal
func TestGetGenres_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := context.Background()

	genre := repo.Genre{
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	id, err := s.CreateGenre(ctx, genre)
	require.NoError(t, err)

	genres, err := s.GetGenres(ctx)
	require.NoError(t, err)
	require.Equal(t, len(genres), 1, "genres len should be 1")

	want := genre
	got := genres[0]
	require.Equal(t, id, got.ID, "id should be equal")
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}
