package repo_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/require"
)

// TestAddRating_Success_ShouldBeNoError tests case when we add user rating, and there should be no error
func TestAddRating_Success_ShouldBeNoError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()
	gameID, err := s.CreateGame(context.Background(), cg)
	require.NoError(t, err, "err should be nil")

	rating := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: td.String(),
		GameID: gameID,
	}

	err = s.AddRating(context.Background(), rating)
	require.NoError(t, err, "err should be nil")
}

// TestGetUserRatingsByGamesIDs_DataExists_ShouldBeEqual tests case when we add user rating, then get user ratings for
// specified games, and matching data should be equal
func TestGetUserRatingsByGamesIDs_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg1 := getCreateGameData()
	gameID1, err := s.CreateGame(context.Background(), cg1)
	require.NoError(t, err, "err should be nil")
	cg2 := getCreateGameData()
	gameID2, err := s.CreateGame(context.Background(), cg2)
	require.NoError(t, err, "err should be nil")
	cg3 := getCreateGameData()
	gameID3, err := s.CreateGame(context.Background(), cg3)
	require.NoError(t, err, "err should be nil")

	userID := td.String()
	rating1 := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: userID,
		GameID: gameID1,
	}
	rating2 := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: userID,
		GameID: gameID2,
	}
	rating3 := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: td.String(),
		GameID: gameID3,
	}

	err = s.AddRating(context.Background(), rating1)
	require.NoError(t, err, "err should be nil")
	err = s.AddRating(context.Background(), rating2)
	require.NoError(t, err, "err should be nil")
	err = s.AddRating(context.Background(), rating3)
	require.NoError(t, err, "err should be nil")

	ratings, err := s.GetUserRatingsByGamesIDs(context.Background(), userID, []int32{rating1.GameID, rating3.GameID})
	require.NoError(t, err, "err should be nil")
	require.Equal(t, len(ratings), 1, "ratings len should be 1")

	want := rating1
	got := ratings[0]
	require.Equal(t, want.GameID, got.GameID, "game id should be equal")
	require.Equal(t, want.UserID, got.UserID, "user id should be equal")
	require.Equal(t, want.Rating, got.Rating, "rating should be equal")
}

// TestGetUserRatingsByGamesIDs_DataExists_ShouldBeEqual tests case when we add user rating, then get user ratings for
// specified games, and matching data should be equal
func TestGetUserRatings_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg1 := getCreateGameData()
	gameID1, err := s.CreateGame(context.Background(), cg1)
	require.NoError(t, err, "err should be nil")
	cg2 := getCreateGameData()
	gameID2, err := s.CreateGame(context.Background(), cg2)
	require.NoError(t, err, "err should be nil")
	cg3 := getCreateGameData()
	gameID3, err := s.CreateGame(context.Background(), cg3)
	require.NoError(t, err, "err should be nil")

	userID := td.String()
	rating1 := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: userID,
		GameID: gameID1,
	}
	rating2 := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: userID,
		GameID: gameID2,
	}
	rating3 := repo.CreateRating{
		Rating: td.Uint8(),
		UserID: td.String(),
		GameID: gameID3,
	}

	err = s.AddRating(context.Background(), rating1)
	require.NoError(t, err, "err should be nil")
	err = s.AddRating(context.Background(), rating2)
	require.NoError(t, err, "err should be nil")
	err = s.AddRating(context.Background(), rating3)
	require.NoError(t, err, "err should be nil")

	ratings, err := s.GetUserRatings(context.Background(), userID)
	require.NoError(t, err, "err should be nil")
	require.Equal(t, len(ratings), 2, "ratings len should be 2")

	want1, want2 := rating1, rating2
	got1, ok1 := ratings[rating1.GameID]
	got2, ok2 := ratings[rating2.GameID]
	require.True(t, ok1, "rating should exist")
	require.True(t, ok2, "rating should exist")
	require.Equal(t, want1.Rating, got1, "game id should be equal")
	require.Equal(t, want2.Rating, got2, "game id should be equal")
}
