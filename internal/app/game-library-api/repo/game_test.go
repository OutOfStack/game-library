package repo_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/stretchr/testify/require"
)

// TestGetGames_NotExist_ShouldReturnEmpty tests case when there is no data, and we should get empty result
func TestGetGames_NotExist_ShouldReturnEmpty(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	games, err := s.GetGames(t.Context(), 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault})
	require.NoError(t, err)

	require.Empty(t, games, "games should be empty")
}

// TestGetGames_DataExists_ShouldBeEqual tests case when we add one game, then fetch first game, and they should be equal
func TestGetGames_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()

	_, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	games, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault})
	require.NoError(t, err)

	require.Len(t, games, 1, "len of games should be 1")

	want := cg
	got := games[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGames_OrderByDefault_ShouldReturnOrdered tests case when we add two game, then order by default and sort should work correctly
func TestGetGames_OrderByDefault_ShouldReturnOrdered(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg1 := getCreateGameData()
	cg1.ReleaseDate = "1990-01-01"
	cg1.IGDBRating = 100
	cg2 := getCreateGameData()
	cg2.ReleaseDate = "2024-01-01"
	cg2.IGDBRating = 95

	_, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	games, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault})
	require.NoError(t, err)

	require.Len(t, games, 2, "len of games should be 2")

	// weight of game 2 is higher
	want := cg2
	got := games[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGames_OrderByName_ShouldReturnOrdered tests case when we add two game, then order by name and sort should work correctly
func TestGetGames_OrderByName_ShouldReturnOrdered(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg1 := getCreateGameData()
	cg1.Name = "Yakuza"
	cg2 := getCreateGameData()
	cg2.Name = "Adventure Time"

	_, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	games, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByName})
	require.NoError(t, err)

	require.Len(t, games, 2, "len of games should be 2")

	want := cg2
	got := games[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGames_OrderByReleaseDate_ShouldReturnOrdered tests case when we add two game, then order by release date and sort should work correctly
func TestGetGames_OrderByReleaseDate_ShouldReturnOrdered(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg1 := getCreateGameData()
	cg1.ReleaseDate = "1991-01-01"
	cg2 := getCreateGameData()
	cg2.ReleaseDate = "2023-01-01"

	_, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	games, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByReleaseDate})
	require.NoError(t, err)

	require.Len(t, games, 2, "len of games should be 2")

	want := cg2
	got := games[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGames_FilterByName_ShouldReturnEqual tests case when we add game, then filter this game by name, and they should be equal
func TestGetGames_FilterByName_ShouldReturnEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()

	_, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	matched, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault, Name: cg.Name})
	require.NoError(t, err)

	require.Len(t, matched, 1, "len of matched should be 1")

	want := cg
	got := matched[0]
	compareCreateGameAndGame(t, want, got)
}

// TestGetGames_FilterByName_ShouldReturnMatched tests case when we add multiple games, then filter games by name, and we should get matches
func TestGetGames_FilterByName_ShouldReturnMatched(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	ng1 := getCreateGameData()
	ng1.Name = "test game name"
	ng2 := getCreateGameData()
	ng2.Name = "tEsTGameName"
	ng3 := getCreateGameData()
	ng3.Name = "tEssTGameName"
	ng4 := getCreateGameData()
	ng4.Name = "a TEST game name"

	_, err := s.CreateGame(ctx, ng1)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng2)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng3)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng4)
	require.NoError(t, err)

	matched, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault, Name: "test"})
	require.NoError(t, err)

	// ng1, ng2, ng4
	require.Len(t, matched, 3, "len should be 3")
}

// TestGetGames_Filte_ShouldReturnMatched tests case when we add multiple games, then filter games by developer, publisher and genre and we should get matches
func TestGetGames_Filter_ShouldReturnMatched(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	developer1, developer2 := td.Int32(), td.Int32()
	genre1, genre2 := td.Int32(), td.Int32()
	publisher1, publisher2 := td.Int32(), td.Int32()

	ng1 := getCreateGameData()
	ng1.DevelopersIDs = append(ng1.DevelopersIDs, developer1, developer2)
	ng1.GenresIDs = append(ng1.GenresIDs, genre1, genre2)
	ng1.PublishersIDs = append(ng1.PublishersIDs, publisher1)
	ng2 := getCreateGameData()
	ng2.DevelopersIDs = append(ng2.DevelopersIDs, developer1, developer2)
	ng2.GenresIDs = append(ng2.GenresIDs, genre2)
	ng2.PublishersIDs = append(ng2.PublishersIDs, publisher2)
	ng3 := getCreateGameData()
	ng3.DevelopersIDs = append(ng3.DevelopersIDs, developer1, developer2)
	ng3.GenresIDs = append(ng3.GenresIDs, genre1)
	ng3.PublishersIDs = append(ng3.PublishersIDs, publisher1, publisher2)

	_, err := s.CreateGame(ctx, ng1)
	require.NoError(t, err)
	g2ID, err := s.CreateGame(ctx, ng2)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng3)
	require.NoError(t, err)

	matched, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault, DeveloperID: developer2, PublisherID: publisher2, GenreID: genre2})
	require.NoError(t, err)

	// ng2
	require.Len(t, matched, 1, "len should be 1")
	require.Equal(t, g2ID, matched[0].ID, "games ids should match")
}

// TestGetGamesCount_DataExists_ShouldReturnCount tests case when we add multiple games, get their count, and it should match
func TestGetGamesCount_DataExists_ShouldReturnCount(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	ng1 := getCreateGameData()
	ng2 := getCreateGameData()

	_, err := s.CreateGame(ctx, ng1)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng2)
	require.NoError(t, err)

	count, err := s.GetGamesCount(ctx, model.GamesFilter{})
	require.NoError(t, err)

	require.Equal(t, 2, int(count), "count should be 2")
}

// TestGetGamesCount_FilterByName_ShouldReturnMatchedCount tests case when we add multiple games, get count by name, amd we should get count of matches
func TestGetGamesCount_FilterByName_ShouldReturnMatchedCount(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	ng1 := getCreateGameData()
	ng1.Name = "the best Game"
	ng2 := getCreateGameData()
	ng2.Name = "game, The best"
	ng3 := getCreateGameData()
	ng3.Name = "GAME"
	ng4 := getCreateGameData()
	ng4.Name = "the Gane"

	_, err := s.CreateGame(ctx, ng1)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng2)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng3)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng4)
	require.NoError(t, err)

	count, err := s.GetGamesCount(ctx, model.GamesFilter{Name: "game"})
	require.NoError(t, err)

	// ng1, ng2, ng3
	require.Equal(t, 3, int(count), "count should be 3")
}

// TestGetGamesCount_Filter_ShouldReturnMatched tests case when we add multiple games, get count by developer, publisher and genre and we should get count of matches
func TestGetGamesCount_Filter_ShouldReturnMatched(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	developer1, developer2 := td.Int32(), td.Int32()
	genre1, genre2 := td.Int32(), td.Int32()
	publisher1, publisher2 := td.Int32(), td.Int32()

	ng1 := getCreateGameData()
	ng1.DevelopersIDs = append(ng1.DevelopersIDs, developer1, developer2)
	ng1.GenresIDs = append(ng1.GenresIDs, genre1, genre2)
	ng1.PublishersIDs = append(ng1.PublishersIDs, publisher1)
	ng2 := getCreateGameData()
	ng2.DevelopersIDs = append(ng2.DevelopersIDs, developer1, developer2)
	ng2.GenresIDs = append(ng2.GenresIDs, genre2)
	ng2.PublishersIDs = append(ng2.PublishersIDs, publisher2)
	ng3 := getCreateGameData()
	ng3.DevelopersIDs = append(ng3.DevelopersIDs, developer1, developer2)
	ng3.GenresIDs = append(ng3.GenresIDs, genre1)
	ng3.PublishersIDs = append(ng3.PublishersIDs, publisher1, publisher2)

	_, err := s.CreateGame(ctx, ng1)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng2)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, ng3)
	require.NoError(t, err)

	count, err := s.GetGamesCount(ctx, model.GamesFilter{OrderBy: repo.OrderGamesByDefault, DeveloperID: developer1, PublisherID: publisher1, GenreID: genre1})
	require.NoError(t, err)

	// ng1, ng3
	require.EqualValues(t, 2, count, "count should be 2")
}

// TestGetGamesCount_DataNotExist_ShouldReturnZero tests case when there ara no games, get their count, and it should be 0
func TestGetGamesCount_DataNotExist_ShouldReturnZero(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	count, err := s.GetGamesCount(t.Context(), model.GamesFilter{Name: ""})
	require.NoError(t, err)

	require.Equal(t, 0, int(count), "len of matched should be 0")
}

// TestGetGameByID_NotExist_ShouldReturnNotFoundError tests case when a game with provided id does not exist, and we should get a Not Found Error
func TestGetGameByID_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	g, err := s.GetGameByID(t.Context(), id)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", id), "err should be NotFound")
	require.Zero(t, g.ID, "id should be 0")
}

// TestGetGameByID_DataExists_ShouldRetrieveEqual tests case when we add game, then fetch this game, and they should be equal
func TestGetGameByID_DataExists_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()

	id, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	gg, err := s.GetGameByID(ctx, id)
	require.NoError(t, err)

	want := cg
	got := gg
	compareCreateGameAndGame(t, want, got)
}

// TestUpdateGame_Valid_ShouldRetrieveEqual tests case when we update game, then fetch this game, and they should be equal
func TestUpdateGame_Valid_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cr := getCreateGameData()

	id, err := s.CreateGame(ctx, cr)
	require.NoError(t, err)

	up := model.UpdateGameData{
		Name:         td.String(),
		Developers:   []int32{td.Int32(), td.Int32()},
		Publishers:   []int32{td.Int32(), td.Int32()},
		ReleaseDate:  types.DateOf(td.Date()).String(),
		Genres:       []int32{td.Int32(), td.Int32()},
		LogoURL:      td.String(),
		Summary:      td.String(),
		Slug:         td.String(),
		PlatformsIDs: []int32{td.Int32(), td.Int32()},
		Screenshots:  []string{td.String(), td.String()},
		Websites:     []string{td.String(), td.String()},
		IGDBRating:   td.Float64(),
		IGDBID:       int64(td.Uint32()),
	}

	err = s.UpdateGame(ctx, id, up)
	require.NoError(t, err)

	g, err := s.GetGameByID(ctx, id)
	require.NoError(t, err)

	want := up
	got := g
	compareUpdateGameAndGame(t, want, got)
}

// TestUpdateGame_NotExist_ShouldReturnNotFoundError tests case when we update a non-existing game, and we should get a Not Found Error
func TestUpdateGame_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	up := model.UpdateGameData{ReleaseDate: "2022-05-18"}
	err := s.UpdateGame(t.Context(), id, up)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", id), "err should be NotFound")
}

// TestDeleteGame_Valid_ShouldDelete tests case when we delete a game
func TestDeleteGame_Valid_ShouldDelete(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cr := getCreateGameData()
	id, err := s.CreateGame(ctx, cr)
	require.NoError(t, err)

	err = s.DeleteGame(ctx, id)
	require.NoError(t, err)

	g, err := s.GetGameByID(ctx, id)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", id), "err should be NotFound")
	require.Zero(t, g.ID, "id should be 0")
}

// TestUpdateRating_Valid_ShouldUpdateGameRating tests case when we update game rating
func TestUpdateRating_Valid_ShouldUpdateGameRating(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cr := getCreateGameData()
	id, err := s.CreateGame(ctx, cr)
	require.NoError(t, err)

	var r1, r2, r3 = td.Uint8(), td.Uint8(), td.Uint8()
	err = s.AddRating(ctx, model.CreateRating{Rating: r1, UserID: td.String(), GameID: id})
	require.NoError(t, err)
	err = s.AddRating(ctx, model.CreateRating{Rating: r2, UserID: td.String(), GameID: id})
	require.NoError(t, err)
	err = s.AddRating(ctx, model.CreateRating{Rating: r3, UserID: td.String(), GameID: id})
	require.NoError(t, err)

	err = s.UpdateGameRating(ctx, id)
	require.NoError(t, err)

	game, err := s.GetGameByID(ctx, id)
	require.NoError(t, err)

	sum := int(r1) + int(r2) + int(r3)
	want := float64(sum) / 3
	got := game.Rating
	require.InDelta(t, want, got, 0.01, "rating should be in delta 0.01")
}

func getCreateGameData() model.CreateGame {
	return model.CreateGame{
		Name:          td.String(),
		DevelopersIDs: []int32{td.Int32(), td.Int32()},
		PublishersIDs: []int32{td.Int32(), td.Int32()},
		ReleaseDate:   td.Date().Format("2006-01-02"),
		GenresIDs:     []int32{td.Int32(), td.Int32()},
		LogoURL:       td.String(),
		Summary:       td.String(),
		Slug:          td.String(),
		PlatformsIDs:  []int32{td.Int32(), td.Int32()},
		Screenshots:   []string{td.String(), td.String()},
		Websites:      []string{td.String(), td.String()},
		IGDBRating:    td.Float64(),
		IGDBID:        int64(td.Uint32()),
	}
}

func compareCreateGameAndGame(t *testing.T, want model.CreateGame, got model.Game) {
	t.Helper()

	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.DevelopersIDs, got.DevelopersIDs, "developers should be equal")
	require.Equal(t, want.PublishersIDs, got.PublishersIDs, "publisher should be equal")
	require.Equal(t, want.ReleaseDate, got.ReleaseDate.String(), "release date should be equal")
	require.Equal(t, want.GenresIDs, got.GenresIDs, "genres should be equal")
	require.Equal(t, want.LogoURL, got.LogoURL, "logo url should be equal")
	require.Equal(t, want.Summary, got.Summary, "summary should be equal")
	require.Equal(t, want.Slug, got.Slug, "slug should be equal")
	require.Equal(t, want.PlatformsIDs, got.PlatformsIDs, "platforms should be equal")
	require.Equal(t, want.Screenshots, got.Screenshots, "screenshots should be equal")
	require.Equal(t, want.Websites, got.Websites, "websites should be equal")
	require.InDeltaf(t, want.IGDBRating, got.IGDBRating, 0.01, "igdb rating should be almost equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

func compareUpdateGameAndGame(t *testing.T, want model.UpdateGameData, got model.Game) {
	t.Helper()

	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.Developers, got.DevelopersIDs, "developers should be equal")
	require.Equal(t, want.Publishers, got.PublishersIDs, "publisher should be equal")
	require.Equal(t, want.ReleaseDate, got.ReleaseDate.String(), "release date should be equal")
	require.Equal(t, want.Genres, got.GenresIDs, "genres should be equal")
	require.Equal(t, want.LogoURL, got.LogoURL, "logo url should be equal")
	require.Equal(t, want.Summary, got.Summary, "summary should be equal")
	require.Equal(t, want.Slug, got.Slug, "slug should be equal")
	require.Equal(t, want.PlatformsIDs, got.PlatformsIDs, "platforms should be equal")
	require.Equal(t, want.Screenshots, got.Screenshots, "screenshots should be equal")
	require.Equal(t, want.Websites, got.Websites, "websites should be equal")
	require.InDeltaf(t, want.IGDBRating, got.IGDBRating, 0.01, "igdb rating should be almost equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}
