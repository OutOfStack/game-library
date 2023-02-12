package repo_test

import (
	"context"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/stretchr/testify/require"
)

// TestGetGames_NotExist_ShouldReturnEmpty tests case when there is no data, and we should get empty result
func TestGetGames_NotExist_ShouldReturnEmpty(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	games, err := s.GetGames(context.Background(), 20, 0)
	require.NoError(t, err, "err should be nil")

	require.Zero(t, len(games), "len of games should be 0")
}

// TestGetGames_DataExists_ShouldBeEqual tests case when we add one game, then fetch first game, and they should be equal
func TestGetGames_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()

	_, err := s.CreateGame(context.Background(), cg)
	require.NoError(t, err, "err should be nil")

	games, err := s.GetGames(context.Background(), 20, 0)
	require.NoError(t, err, "err should be nil")

	require.Equal(t, 1, len(games), "len of games should be 1")

	want := cg
	got := games[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGameByID_NotExist_ShouldReturnNotFoundError tests case when a game with provided id does not exist, and we should get a Not Found Error
func TestGetGameByID_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	g, err := s.GetGameByID(context.Background(), id)
	require.ErrorIs(t, err, repo.ErrNotFound[int32]{Entity: "game", ID: id}, "err should be NotFound")
	require.Zero(t, g.ID, "id should be 0")
}

// TestGetGameByID_DataExists_ShouldRetrieveEqual tests case when we add game, then fetch this game, and they should be equal
func TestGetGameByID_DataExists_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()

	id, err := s.CreateGame(context.Background(), cg)
	require.NoError(t, err, "err should be nil")

	gg, err := s.GetGameByID(context.Background(), id)
	require.NoError(t, err, "err should be nil")

	want := cg
	got := gg
	compareCreateGameAndGame(t, want, got)
}

// TestSearchGames_DataExists_ShouldReturnEqual tests case when we add game, then search this game, and they should be equal
func TestSearchGames_DataExists_ShouldReturnEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cg := getCreateGameData()

	_, err := s.CreateGame(context.Background(), cg)
	require.NoError(t, err, "err should be nil")

	matched, err := s.SearchGames(context.Background(), cg.Name)
	require.NoError(t, err, "err should be nil")

	require.Equal(t, 1, len(matched), "len of matched should be 1")

	want := cg
	got := matched[0]
	compareCreateGameAndGame(t, want, got)
}

// TestSearchGames_DataExists_ShouldReturnMatched tests case when we add multiple games, then search games, and we should get matches
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
	require.NoError(t, err, "err should be nil")
	_, err = s.CreateGame(context.Background(), ng2)
	require.NoError(t, err, "err should be nil")
	_, err = s.CreateGame(context.Background(), ng3)
	require.NoError(t, err, "err should be nil")
	_, err = s.CreateGame(context.Background(), ng4)
	require.NoError(t, err, "err should be nil")

	matched, err := s.SearchGames(context.Background(), "test")
	require.NoError(t, err, "err should be nil")

	require.Equal(t, 2, len(matched), "len of matched should be 2")
}

// TestUpdateGame_Valid_ShouldRetrieveEqual tests case when we update game, then fetch this game, and they should be equal
func TestUpdateGame_Valid_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()

	id, err := s.CreateGame(context.Background(), cr)
	require.NoError(t, err, "err should be nil")

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
	require.NoError(t, err, "err should be nil")

	g, err := s.GetGameByID(context.Background(), id)
	require.NoError(t, err, "err should be nil")

	want := up
	got := g
	compareUpdateGameAndGame(t, want, got)
}

// TestUpdateGame_NotExist_ShouldReturnNotFoundError tests case when we update a non-existing game, and we should get a Not Found Error
func TestUpdateGame_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	up := repo.UpdateGame{ReleaseDate: "2022-05-18"}
	err := s.UpdateGame(context.Background(), id, up)
	require.ErrorIs(t, err, repo.ErrNotFound[int32]{Entity: "game", ID: id}, "err should be NotFound")
}

// TestDeleteGame_Valid_ShouldDelete tests case when we delete a game
func TestDeleteGame_Valid_ShouldDelete(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()
	id, err := s.CreateGame(context.Background(), cr)
	require.NoError(t, err, "err should be nil")

	err = s.DeleteGame(context.Background(), id)
	require.NoError(t, err, "err should be nil")

	g, err := s.GetGameByID(context.Background(), id)
	require.ErrorIs(t, err, repo.ErrNotFound[int32]{Entity: "game", ID: id}, "err should be NotFound")
	require.Zero(t, g.ID, "id should be 0")
}

// TestUpdateRating_Valid_ShouldUpdateGameRating tests case when we update game rating
func TestUpdateRating_Valid_ShouldUpdateGameRating(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()
	id, err := s.CreateGame(context.Background(), cr)
	require.NoError(t, err, "err should be nil")

	var r1, r2, r3 uint8 = td.Uint8(), td.Uint8(), td.Uint8()
	err = s.AddRating(context.Background(), repo.CreateRating{Rating: r1, UserID: td.String(), GameID: id})
	require.NoError(t, err, "err should be nil")
	err = s.AddRating(context.Background(), repo.CreateRating{Rating: r2, UserID: td.String(), GameID: id})
	require.NoError(t, err, "err should be nil")
	err = s.AddRating(context.Background(), repo.CreateRating{Rating: r3, UserID: td.String(), GameID: id})
	require.NoError(t, err, "err should be nil")

	err = s.UpdateGameRating(context.Background(), id)
	require.NoError(t, err, "err should be nil")

	game, err := s.GetGameByID(context.Background(), id)
	require.NoError(t, err, "err should be nil")

	sum := int(r1) + int(r2) + int(r3)
	want := float64(sum) / 3
	got := game.Rating
	require.InDelta(t, want, got, 0.01, "rating should be in delta 0.01")
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
	require.InDeltaf(t, want.IGDBRating, got.IGDBRating, 0.01, "igdb rating should be almost equal")
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
	require.InDeltaf(t, want.IGDBRating, got.IGDBRating, 0.01, "igdb rating should be almost equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}
