package repo_test

import (
	"context"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

// TestCreateGenre_IGDBIDIsNull_ShouldBeNoError tests case when we add genre without igdb id, and there should be no error
func TestCreateGenre_IGDBIDIsNull_ShouldBeNoError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	genre := model.Genre{
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

	genre := model.Genre{
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

func TestGetTopGenres_Ok(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := context.Background()

	// create 2 genres and 3 games
	genre1ID, err := s.CreateGenre(ctx, model.Genre{
		Name:   td.String(),
		IGDBID: td.Int64(),
	})
	require.NoError(t, err)

	genre2ID, err := s.CreateGenre(ctx, model.Genre{
		Name:   td.String(),
		IGDBID: td.Int64(),
	})
	require.NoError(t, err)

	cg1, cg2, cg3 := getCreateGameData(), getCreateGameData(), getCreateGameData()

	// genre 1 is in 2 games, genre 2 is in 3 games
	cg1.Genres = []int32{genre1ID, genre2ID}
	cg2.Genres = []int32{genre2ID}
	cg3.Genres = []int32{genre1ID, genre2ID}

	_, err = s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	top, err := s.GetTopGenres(ctx, 5)
	require.NoError(t, err)

	require.Len(t, top, 2, "len of top genres should be 3")

	require.Equal(t, genre2ID, top[0].ID, "top 1 genre should be genre 2")
	require.Equal(t, genre1ID, top[1].ID, "top 2 genre should be genre 1")
}
