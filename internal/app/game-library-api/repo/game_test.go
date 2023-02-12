package repo_test

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/stretchr/testify/require"
)

// TestGetGames_NotExist_ShouldReturnEmpty tests case when there is no data and we should get empty result
func TestGetGames_NotExist_ShouldReturnEmpty(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	games, err := s.GetGames(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("error getting games: %v", err)
	}

	want := 0
	got := len(games)
	if want != got {
		t.Fatalf("Expected to retrieve %d games, got %d", want, got)
	}
}

// TestGetGames_DataExists_ShouldBeEqual tests case when we add one game, then fetch first game and they should be equal
func TestGetGames_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()

	_, err := s.CreateGame(context.Background(), cg)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	games, err := s.GetGames(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("error getting games: %v", err)
	}

	if len(games) != 1 {
		t.Fatalf("Expected to retrieve 1 game, got %d", len(games))
	}

	want := cg
	got := games[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGameByID_NotExist_ShouldReturnNotFoundError tests case when a game with provided id does not exist and we should get a Not Found Error
func TestGetGameByID_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	g, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		if !errors.As(err, &repo.ErrNotFound[int32]{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = repo.ErrNotFound[int32]{
		Entity: "game",
		ID:     id,
	}
	gotErr := err
	if g.ID != 0 || !errors.Is(gotErr, wantErr) {
		t.Fatalf("Expected to receive empty entity and error [%v], got [%v]", wantErr, gotErr)
	}
}

// TestGetGameByID_DataExists_ShouldRetrieveEqual tests case when we add game, then fetch this game and they should be equal
func TestGetGameByID_DataExists_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()

	id, err := s.CreateGame(context.Background(), cg)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	gg, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	want := cg
	got := gg
	compareCreateGameAndGame(t, want, got)
}

// TestSearchGames_DataExists_ShouldReturnEqual tests case when we add game, then search this game and they should be equal
func TestSearchGames_DataExists_ShouldReturnEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()

	_, err := s.CreateGame(context.Background(), cg)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	matched, err := s.SearchGames(context.Background(), cg.Name)
	if err != nil {
		t.Fatalf("error searching games: %v", err)
	}

	if len(matched) != 1 {
		t.Fatalf("Expected to retrieve 1 matched game, got %d", len(matched))
	}

	want := cg
	got := matched[0]
	compareCreateGameAndGame(t, want, got)
}

// TestSearchGames_DataExists_ShouldReturnMatched tests case when we add multiple games, then search games and we should get matches
func TestSearchInfos_DataExists_ShouldReturnMatched(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ng1 := getCreateGameData()
	ng1.Name = "test game name"
	ng2 := getCreateGameData()
	ng2.Name = "tEsTGameName"
	ng3 := getCreateGameData()
	ng3.Name = "tEssTGameName"
	ng4 := getCreateGameData()
	ng4.Name = "a test game name"

	_, err := s.CreateGame(context.Background(), ng1)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}
	_, err = s.CreateGame(context.Background(), ng2)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}
	_, err = s.CreateGame(context.Background(), ng3)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}
	_, err = s.CreateGame(context.Background(), ng4)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	matched, err := s.SearchGames(context.Background(), "test")
	if err != nil {
		t.Fatalf("error searching games: %v", err)
	}

	want := 2
	got := len(matched)
	if want != got {
		t.Fatalf("Expected to retrieve %d matched game, got %d", want, got)
	}
}

// TestUpdateGame_Valid_ShouldRetrieveEqual tests case when we update game, then fetch this game and they should be equal
func TestUpdateGame_Valid_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()

	id, err := s.CreateGame(context.Background(), cr)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	up := repo.UpdateGame{
		Name:        td.String(),
		Developers:  []int32{td.Int32(), td.Int32()},
		Publishers:  []int32{td.Int32(), td.Int32()},
		ReleaseDate: types.DateOf(td.Date()).String(),
		Genres:      []int32{td.Int32(), td.Int32()},
		LogoURL:     td.String(),
		Summary:     td.String(),
		Slug:        td.String(),
		Platforms:   []int32{td.Int32(), td.Int32()},
		Screenshots: []string{td.String(), td.String()},
		Websites:    []string{td.String(), td.String()},
		IGDBRating:  td.Float64(),
		IGDBID:      int64(td.Uint32()),
	}

	err = s.UpdateGame(context.Background(), id, up)
	if err != nil {
		t.Fatalf("error updating game: %v", err)
	}

	g, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	want := up
	got := g
	compareUpdateGameAndGame(t, want, got)
}

// TestUpdateGame_NotExist_ShouldReturnNotFoundError tests case when we update a non existing game and we should get a Not Found Error
func TestUpdateGame_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	up := repo.UpdateGame{ReleaseDate: "2022-05-18"}
	err := s.UpdateGame(context.Background(), id, up)
	if err != nil {
		if !errors.As(err, &repo.ErrNotFound[int32]{}) {
			t.Fatalf("error updating game: %v", err)
		}
	}

	var wantErr = repo.ErrNotFound[int32]{
		Entity: "game",
		ID:     id,
	}
	gotErr := err
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("Expected to get empty entity and error [%v], got [%v]", wantErr, gotErr)
	}
}

// TestDeleteGame_Valid_ShouldDelete tests case when we delete a game
func TestDeleteGame_Valid_ShouldDelete(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()
	id, err := s.CreateGame(context.Background(), cr)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	err = s.DeleteGame(context.Background(), id)
	if err != nil {
		t.Fatalf("error deleting game: %v", err)
	}

	g, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		if !errors.As(err, &repo.ErrNotFound[int32]{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = repo.ErrNotFound[int32]{
		Entity: "game",
		ID:     id,
	}
	gotErr := err
	if g.ID != 0 || !errors.Is(gotErr, wantErr) {
		t.Fatalf("Expected to receive empty entity and error [%v], got [%v]", wantErr, gotErr)
	}
}

// TestUpdateRating_Valid_ShouldUpdateGameRating tests case when we update game rating
func TestUpdateRating_Valid_ShouldUpdateGameRating(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()
	id, err := s.CreateGame(context.Background(), cr)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	var r1, r2, r3 uint8 = td.Uint8(), td.Uint8(), td.Uint8()
	err = s.AddRating(context.Background(), repo.CreateRating{Rating: r1, UserID: td.String(), GameID: id})
	if err != nil {
		t.Fatalf("error adding rating: %v", err)
	}
	err = s.AddRating(context.Background(), repo.CreateRating{Rating: r2, UserID: td.String(), GameID: id})
	if err != nil {
		t.Fatalf("error adding rating: %v", err)
	}
	err = s.AddRating(context.Background(), repo.CreateRating{Rating: r3, UserID: td.String(), GameID: id})
	if err != nil {
		t.Fatalf("error adding rating: %v", err)
	}

	err = s.UpdateGameRating(context.Background(), id)
	if err != nil {
		if !errors.As(err, &repo.ErrNotFound[int32]{}) {
			t.Fatalf("error updating game: %v", err)
		}
	}

	game, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	sum := int(r1) + int(r2) + int(r3)
	want := float64(sum) / 3
	got := game.Rating
	if int(want) != int(got) {
		t.Errorf("Expected to get game rating with value %f, got %f", want, got)
	}

	wantDelta := 0.01
	gotDelta := math.Abs(want - got)
	if gotDelta > wantDelta {
		t.Errorf("Expected delta to be nor more than %f, got %f", wantDelta, gotDelta)
	}
}

func getCreateGameData() repo.CreateGame {
	return repo.CreateGame{
		Name:        td.String(),
		Developer:   td.String(),
		Developers:  []int32{td.Int32(), td.Int32()},
		Publisher:   td.String(),
		Publishers:  []int32{td.Int32(), td.Int32()},
		ReleaseDate: td.Date().Format("2006-01-02"),
		Genre:       []string{td.String(), td.String()},
		Genres:      []int32{td.Int32(), td.Int32()},
		LogoURL:     td.String(),
		Summary:     td.String(),
		Slug:        td.String(),
		Platforms:   []int32{td.Int32(), td.Int32()},
		Screenshots: []string{td.String(), td.String()},
		Websites:    []string{td.String(), td.String()},
		IGDBRating:  td.Float64(),
		IGDBID:      int64(td.Uint32()),
	}
}

func compareCreateGameAndGame(t *testing.T, want repo.CreateGame, got repo.Game) {
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.Developer, got.Developer, "developer should be equal")
	require.Equal(t, want.Developers, []int32(got.Developers), "developers should be equal")
	require.Equal(t, want.Publisher, got.Publisher, "publisher should be equal")
	require.Equal(t, want.Publishers, []int32(got.Publishers), "publisher should be equal")
	require.Equal(t, want.ReleaseDate, got.ReleaseDate.String(), "release date should be equal")
	require.Equal(t, want.Genre, []string(got.Genre), "genre should be equal")
	require.Equal(t, want.Genres, []int32(got.Genres), "genres should be equal")
	require.Equal(t, want.LogoURL, got.LogoURL, "logo url should be equal")
	require.Equal(t, want.Summary, got.Summary, "summary should be equal")
	require.Equal(t, want.Slug, got.Slug, "slug should be equal")
	require.Equal(t, want.Platforms, []int32(got.Platforms), "platforms should be equal")
	require.Equal(t, want.Screenshots, []string(got.Screenshots), "screenshots should be equal")
	require.Equal(t, want.Websites, []string(got.Websites), "websites should be equal")
	require.InDeltaf(t, want.IGDBRating, got.IGDBRating, 0.1, "igdb rating should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

func compareUpdateGameAndGame(t *testing.T, want repo.UpdateGame, got repo.Game) {
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.Developers, []int32(got.Developers), "developers should be equal")
	require.Equal(t, want.Publishers, []int32(got.Publishers), "publisher should be equal")
	require.Equal(t, want.ReleaseDate, got.ReleaseDate.String(), "release date should be equal")
	require.Equal(t, want.Genres, []int32(got.Genres), "genres should be equal")
	require.Equal(t, want.LogoURL, got.LogoURL, "logo url should be equal")
	require.Equal(t, want.Summary, got.Summary, "summary should be equal")
	require.Equal(t, want.Slug, got.Slug, "slug should be equal")
	require.Equal(t, want.Platforms, []int32(got.Platforms), "platforms should be equal")
	require.Equal(t, want.Screenshots, []string(got.Screenshots), "screenshots should be equal")
	require.Equal(t, want.Websites, []string(got.Websites), "websites should be equal")
	require.InDeltaf(t, want.IGDBRating, got.IGDBRating, 0.1, "igdb rating should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}
