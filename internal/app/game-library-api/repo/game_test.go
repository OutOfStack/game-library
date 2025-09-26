package repo_test

import (
	"testing"
	"time"

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

// TestGetGames_DataExists_ShouldBeEqual tests case when we add one game, then fetch first game and they should be equal
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

	id1, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	id2, err := s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	// Set trending index manually for test
	err = s.UpdateGameTrendingIndex(ctx, id1, 0.5) // Lower index for older game
	require.NoError(t, err)

	err = s.UpdateGameTrendingIndex(ctx, id2, 0.8) // Higher index for newer game
	require.NoError(t, err)

	games, err := s.GetGames(ctx, 20, 1, model.GamesFilter{OrderBy: repo.OrderGamesByDefault})
	require.NoError(t, err)

	require.Len(t, games, 2, "len of games should be 2")

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

// TestGetGames_Filte_ShouldReturnMatched tests case when we add multiple games, then filter games by developer, publisher, and genre, and we should get matches
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

// TestGetGamesCount_Filter_ShouldReturnMatched tests case when we add multiple games, get count by developer, publisher, and genre and we should get count of matches
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
		Name:          td.String(),
		DevelopersIDs: []int32{td.Int32(), td.Int32()},
		PublishersIDs: []int32{td.Int32(), td.Int32()},
		ReleaseDate:   types.DateOf(td.Date()).String(),
		GenresIDs:     []int32{td.Int32(), td.Int32()},
		LogoURL:       td.String(),
		Summary:       td.String(),
		Slug:          td.String(),
		PlatformsIDs:  []int32{td.Int32(), td.Int32()},
		Screenshots:   []string{td.String(), td.String()},
		Websites:      []string{td.String(), td.String()},
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

// TestGetPublisherGamesCount_NoGames_ShouldReturnZero tests case when there are no games in the date range
func TestGetPublisherGamesCount_NoGames_ShouldReturnZero(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	publisherID := td.Int32()

	// set date range to current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	count, err := s.GetPublisherGamesCount(ctx, publisherID, startOfMonth, endOfMonth)
	require.NoError(t, err)
	require.Equal(t, 0, count, "count should be 0")
}

// TestGetPublisherGamesCount_WithGames_ShouldReturnCount tests case when there are games in the date range
func TestGetPublisherGamesCount_WithGames_ShouldReturnCount(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	publisherID := td.Int32()

	// set date range to current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// create games for the publisher
	cg1 := getCreateGameData()
	cg1.PublishersIDs = []int32{publisherID}

	cg2 := getCreateGameData()
	cg2.PublishersIDs = []int32{publisherID}

	// create a game for a different publisher
	cg3 := getCreateGameData()
	cg3.PublishersIDs = []int32{td.Int32()}

	_, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	count, err := s.GetPublisherGamesCount(ctx, publisherID, startOfMonth, endOfMonth)
	require.NoError(t, err)
	require.Equal(t, 2, count, "count should be 2")
}

// TestUpdateGameIGDBInfo_Valid_ShouldUpdateIGDBInfo tests case when we update game IGDB info
func TestUpdateGameIGDBInfo_Valid_ShouldUpdateIGDBInfo(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cr := getCreateGameData()
	id, err := s.CreateGame(ctx, cr)
	require.NoError(t, err)

	igdbData := model.UpdateGameIGDBData{
		Name:            td.String(),
		PlatformsIDs:    []int32{td.Int32(), td.Int32()},
		Websites:        []string{td.String(), td.String()},
		IGDBRating:      td.Float64n(100),
		IGDBRatingCount: int32(td.Intn(10000)),
	}

	err = s.UpdateGameIGDBInfo(ctx, id, igdbData)
	require.NoError(t, err)

	game, err := s.GetGameByID(ctx, id)
	require.NoError(t, err)

	require.Equal(t, igdbData.Name, game.Name, "name should be equal")
	require.Equal(t, igdbData.PlatformsIDs, game.PlatformsIDs, "platforms should be equal")
	require.Equal(t, igdbData.Websites, game.Websites, "websites should be equal")
	require.InDelta(t, igdbData.IGDBRating, game.IGDBRating, 0.01, "igdb rating should be equal")
	require.Equal(t, igdbData.IGDBRatingCount, game.IGDBRatingCount, "igdb rating count should be equal")
}

// TestUpdateGameIGDBInfo_NotExist_ShouldReturnNotFoundError tests case when we update IGDB info of a non-existing game
func TestUpdateGameIGDBInfo_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	igdbData := model.UpdateGameIGDBData{
		Name:            td.String(),
		IGDBRating:      td.Float64n(100),
		IGDBRatingCount: int32(td.Intn(10000)),
	}
	err := s.UpdateGameIGDBInfo(t.Context(), id, igdbData)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", id), "err should be NotFound")
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

func getCreateGameData() model.CreateGameData {
	return model.CreateGameData{
		Name:             td.String(),
		DevelopersIDs:    []int32{td.Int32(), td.Int32()},
		PublishersIDs:    []int32{td.Int32(), td.Int32()},
		ReleaseDate:      td.Date().Format("2006-01-02"),
		GenresIDs:        []int32{td.Int32(), td.Int32()},
		LogoURL:          td.String(),
		Summary:          td.String(),
		Slug:             td.String(),
		PlatformsIDs:     []int32{td.Int32(), td.Int32()},
		Screenshots:      []string{td.String(), td.String()},
		Websites:         []string{td.String(), td.String()},
		IGDBRating:       td.Float64n(100),
		IGDBRatingCount:  int32(td.Intn(10000)),
		IGDBID:           int64(td.Uint32()),
		ModerationStatus: model.ModerationStatusReady,
	}
}

func compareCreateGameAndGame(t *testing.T, want model.CreateGameData, got model.Game) {
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
	require.Equal(t, want.IGDBRatingCount, got.IGDBRatingCount, "igdb rating count should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

func compareUpdateGameAndGame(t *testing.T, want model.UpdateGameData, got model.Game) {
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
}

// TestUpdateGameTrendingIndex_Valid_ShouldUpdateTrendingIndex tests updating trending index
func TestUpdateGameTrendingIndex_Valid_ShouldUpdateTrendingIndex(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()
	id, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	trendingIndex := 0.75
	err = s.UpdateGameTrendingIndex(ctx, id, trendingIndex)
	require.NoError(t, err)

	game, err := s.GetGameByID(ctx, id)
	require.NoError(t, err)

	require.InDelta(t, trendingIndex, game.TrendingIndex, 0.01, "trending index should be equal")
}

// TestUpdateGameTrendingIndex_NotExist_ShouldReturnNotFoundError tests updating non-existing game trending index
func TestUpdateGameTrendingIndex_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	err := s.UpdateGameTrendingIndex(t.Context(), id, 0.5)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", id), "err should be NotFound")
}

// TestGetGameTrendingData_Valid_ShouldReturnData tests getting trending data for a game
func TestGetGameTrendingData_Valid_ShouldReturnData(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create game with specific data
	cg := getCreateGameData()
	cg.ReleaseDate = "2023-06-15"
	cg.IGDBRating = 85.5
	cg.IGDBRatingCount = 999
	id, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	// Add some ratings
	err = s.AddRating(ctx, model.CreateRating{Rating: 4, UserID: td.String(), GameID: id})
	require.NoError(t, err)
	err = s.AddRating(ctx, model.CreateRating{Rating: 5, UserID: td.String(), GameID: id})
	require.NoError(t, err)

	// Update game rating
	err = s.UpdateGameRating(ctx, id)
	require.NoError(t, err)

	data, err := s.GetGameTrendingData(ctx, id)
	require.NoError(t, err)

	require.Equal(t, 2023, data.Year, "year should be 2023")
	require.Equal(t, 6, data.Month, "month should be 6")
	require.InDelta(t, 85.5, data.IGDBRating, 0.01, "IGDB rating should match")
	require.Equal(t, int32(999), data.IGDBRatingCount, "IGDB rating count should match")
	require.InDelta(t, 4.5, data.Rating, 0.01, "user rating should be average of 4 and 5")
	require.Equal(t, int32(2), data.RatingCount, "rating count should be 2")
}

// TestGetGameTrendingData_NotExist_ShouldReturnError tests getting trending data for non-existing game
func TestGetGameTrendingData_NotExist_ShouldReturnError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	_, err := s.GetGameTrendingData(t.Context(), id)
	require.Error(t, err, "should return error for non-existing game")
}

// TestGetGamesIDsAfterID_NoGames_ShouldReturnEmpty tests getting games for trending update when none exist
func TestGetGamesIDsAfterID_NoGames_ShouldReturnEmpty(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	gameIDs, err := s.GetGamesIDsAfterID(t.Context(), 0, 10)
	require.NoError(t, err)
	require.Empty(t, gameIDs, "should return empty slice when no games exist")
}

// TestGetGamesIDsAfterID_WithGames_ShouldReturnOrdered tests getting games for trending update
func TestGetGamesIDsAfterID_WithGames_ShouldReturnOrdered(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create multiple games
	cg1 := getCreateGameData()
	cg2 := getCreateGameData()
	cg3 := getCreateGameData()

	id1, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)
	id2, err := s.CreateGame(ctx, cg2)
	require.NoError(t, err)
	id3, err := s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	// Get all games
	gameIDs, err := s.GetGamesIDsAfterID(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, gameIDs, 3, "should return 3 game IDs")

	// Should be ordered by ID ascending
	expectedIDs := []int32{id1, id2, id3}
	require.Equal(t, expectedIDs, gameIDs, "game IDs should be ordered by ID ascending")
}

// TestGetGamesIDsAfterID_WithOffset_ShouldReturnAfterOffset tests getting games with offset
func TestGetGamesIDsAfterID_WithOffset_ShouldReturnAfterOffset(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create multiple games
	cg1 := getCreateGameData()
	cg2 := getCreateGameData()
	cg3 := getCreateGameData()

	id1, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)
	id2, err := s.CreateGame(ctx, cg2)
	require.NoError(t, err)
	id3, err := s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	// Get games after id1
	gameIDs, err := s.GetGamesIDsAfterID(ctx, id1, 10)
	require.NoError(t, err)
	require.Len(t, gameIDs, 2, "should return 2 game IDs after offset")

	expectedIDs := []int32{id2, id3}
	require.Equal(t, expectedIDs, gameIDs, "should return games after the offset ID")
}

// TestGetGamesIDsAfterID_WithLimit_ShouldRespectLimit tests getting games with batch size limit
func TestGetGamesIDsAfterID_WithLimit_ShouldRespectLimit(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create multiple games
	cg1 := getCreateGameData()
	cg2 := getCreateGameData()
	cg3 := getCreateGameData()

	_, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	// Get only 2 games
	gameIDs, err := s.GetGamesIDsAfterID(ctx, 0, 2)
	require.NoError(t, err)
	require.Len(t, gameIDs, 2, "should respect batch size limit")
}

// TestGetGamesByPublisherID_NoGames_ShouldReturnEmpty tests case when publisher has no games
func TestGetGamesByPublisherID_NoGames_ShouldReturnEmpty(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	publisherID := td.Int32()
	games, err := s.GetGamesByPublisherID(t.Context(), publisherID)
	require.NoError(t, err)
	require.Empty(t, games, "games should be empty")
}

// TestGetGamesByPublisherID_WithGames_ShouldReturnMatches tests case when publisher has games
func TestGetGamesByPublisherID_WithGames_ShouldReturnMatches(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	publisherID := td.Int32()

	// Create games for the publisher
	cg1 := getCreateGameData()
	cg1.PublishersIDs = []int32{publisherID, td.Int32()}

	cg2 := getCreateGameData()
	cg2.PublishersIDs = []int32{td.Int32(), publisherID}

	// Create game for different publisher
	cg3 := getCreateGameData()
	cg3.PublishersIDs = []int32{td.Int32()}

	id1, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)
	id2, err := s.CreateGame(ctx, cg2)
	require.NoError(t, err)
	_, err = s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	games, err := s.GetGamesByPublisherID(ctx, publisherID)
	require.NoError(t, err)
	require.Len(t, games, 2, "should return 2 games for publisher")

	// Games should be ordered by ID descending
	require.Equal(t, id2, games[0].ID, "first game should be the most recent")
	require.Equal(t, id1, games[1].ID, "second game should be the older one")

	// Verify publisher ID is in both games
	require.Contains(t, games[0].PublishersIDs, publisherID, "first game should contain publisher ID")
	require.Contains(t, games[1].PublishersIDs, publisherID, "second game should contain publisher ID")
}

// TestGetGameIDByIGDBID_GameExists_ShouldReturnID tests case when game with IGDB ID exists
func TestGetGameIDByIGDBID_GameExists_ShouldReturnID(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()
	igdbID := int64(td.Uint32())
	cg.IGDBID = igdbID

	expectedID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	actualID, err := s.GetGameIDByIGDBID(ctx, igdbID)
	require.NoError(t, err)
	require.Equal(t, expectedID, actualID, "returned ID should match created game ID")
}

// TestGetGameIDByIGDBID_GameNotExists_ShouldReturnNotFoundError tests case when game with IGDB ID does not exist
func TestGetGameIDByIGDBID_GameNotExists_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	igdbID := int64(td.Uint32())
	id, err := s.GetGameIDByIGDBID(t.Context(), igdbID)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", igdbID), "err should be NotFound")
	require.Zero(t, id, "id should be 0")
}

// TestCreateGame_Valid_ShouldCreateAndReturnID tests case when creating a valid game
func TestCreateGame_Valid_ShouldCreateAndReturnID(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()

	id, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)
	require.NotZero(t, id, "created game ID should not be zero")

	// Verify the game was created correctly
	game, err := s.GetGameByID(ctx, id)
	require.NoError(t, err)
	compareCreateGameAndGame(t, cg, game)
}

// TestUpdateGameModerationID_Valid_ShouldUpdateModerationID tests updating game moderation ID
func TestUpdateGameModerationID_Valid_ShouldUpdateModerationID(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	// Create a moderation record first (required by foreign key constraint)
	createModeration := model.CreateModeration{
		GameID: gameID,
		Status: model.ModerationStatusPending,
		GameData: model.ModerationData{
			Name:      cg.Name,
			Publisher: "Test Publisher",
		},
	}
	moderationID, err := s.CreateModerationRecord(ctx, createModeration)
	require.NoError(t, err)

	err = s.UpdateGameModerationID(ctx, gameID, moderationID)
	require.NoError(t, err)

	game, err := s.GetGameByID(ctx, gameID)
	require.NoError(t, err)
	require.True(t, game.ModerationID.Valid, "moderation ID should be valid")
	require.Equal(t, moderationID, game.ModerationID.Int32, "moderation ID should be updated")
}

// TestUpdateGameModerationID_GameNotExists_ShouldReturnNotFoundError tests updating moderation ID for non-existing game
func TestUpdateGameModerationID_GameNotExists_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	gameID := td.Int32()
	moderationID := td.Int32()
	err := s.UpdateGameModerationID(t.Context(), gameID, moderationID)
	require.ErrorIs(t, err, apperr.NewNotFoundError("game", gameID), "err should be NotFound")
}
