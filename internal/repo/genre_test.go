package repo_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
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

	_, err := s.CreateGenre(t.Context(), genre)
	require.NoError(t, err)
}

// TestGetGenres_DataExists_ShouldBeEqual tests case when we add one genre, then fetch first genre, and they should be equal
func TestGetGenres_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	genre := model.Genre{
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	id, err := s.CreateGenre(ctx, genre)
	require.NoError(t, err)

	genres, err := s.GetGenres(ctx)
	require.NoError(t, err)
	require.Len(t, genres, 1, "genres len should be 1")

	want := genre
	got := genres[0]
	require.Equal(t, id, got.ID, "id should be equal")
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

func TestGetGenreByID_GenreExists_ShouldReturnGenre(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	genre := model.Genre{
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	id, err := s.CreateGenre(ctx, genre)
	require.NoError(t, err)

	gotGenre, err := s.GetGenreByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, gotGenre.ID, "id should be equal")
	require.Equal(t, genre.Name, gotGenre.Name, "name should be equal")
	require.Equal(t, genre.IGDBID, gotGenre.IGDBID, "igdb id should be equal")
}

func TestGetGenreByID_GenreNotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	genre := model.Genre{
		Name:   td.String(),
		IGDBID: td.Int64(),
	}

	_, err := s.CreateGenre(ctx, genre)
	require.NoError(t, err)

	randomID := td.Int32()
	gotGenre, err := s.GetGenreByID(ctx, randomID)
	require.ErrorIs(t, err, apperr.NewNotFoundError("genre", randomID), "err should be NotFound")
	require.Zero(t, gotGenre.ID, "got id should be 0")
}

func TestGetTopGenres_Ok(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

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
	cg1.GenresIDs = []int32{genre1ID, genre2ID}
	cg2.GenresIDs = []int32{genre2ID}
	cg3.GenresIDs = []int32{genre1ID, genre2ID}

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
