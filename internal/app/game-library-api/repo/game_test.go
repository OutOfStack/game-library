package repo_test

import (
	"context"
	"errors"
	"math"
	"testing"

	game "github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
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
		if !errors.As(err, &game.ErrNotFound[int32]{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = game.ErrNotFound[int32]{
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
func TestUpdateGmae_Valid_ShouldRetrieveEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	cr := getCreateGameData()

	id, err := s.CreateGame(context.Background(), cr)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	up := game.UpdateGame{
		Name:        "New game",
		Developer:   "New developer",
		Publisher:   "New publisher",
		ReleaseDate: "2021-11-12",
		Genre:       []string{"adventure"},
		LogoURL:     "https://images/999",
	}

	err = s.UpdateGame(context.Background(), id, up)
	if err != nil {
		t.Fatalf("error updating game: %v", err)
	}

	gg, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	want := up
	got := gg
	compareUpdateGameAndGame(t, want, got)
}

// TestUpdateGame_NotExist_ShouldReturnNotFoundError tests case when we update a non existing game and we should get a Not Found Error
func TestUpdateGame_NotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	up := game.UpdateGame{ReleaseDate: "2022-05-18"}
	err := s.UpdateGame(context.Background(), id, up)
	if err != nil {
		if !errors.As(err, &game.ErrNotFound[int32]{}) {
			t.Fatalf("error updating game: %v", err)
		}
	}

	var wantErr = game.ErrNotFound[int32]{
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
		if !errors.As(err, &game.ErrNotFound[int32]{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = game.ErrNotFound[int32]{
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
	err = s.AddRating(context.Background(), game.CreateRating{Rating: r1, UserID: td.String(), GameID: id})
	if err != nil {
		t.Fatalf("error adding rating: %v", err)
	}
	err = s.AddRating(context.Background(), game.CreateRating{Rating: r2, UserID: td.String(), GameID: id})
	if err != nil {
		t.Fatalf("error adding rating: %v", err)
	}
	err = s.AddRating(context.Background(), game.CreateRating{Rating: r3, UserID: td.String(), GameID: id})
	if err != nil {
		t.Fatalf("error adding rating: %v", err)
	}

	err = s.UpdateGameRating(context.Background(), id)
	if err != nil {
		if !errors.As(err, &game.ErrNotFound[int32]{}) {
			t.Fatalf("error updating game: %v", err)
		}
	}

	game, err := s.GetGameByID(context.Background(), id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	sum := int(r1) + int(r2) + int(r3)
	want := float64(sum) / 3
	got := game.Rating.Float64
	if int(want) != int(got) {
		t.Errorf("Expected to get game rating with value %f, got %f", want, got)
	}

	wantDelta := 0.01
	gotDelta := math.Abs(want - got)
	if gotDelta > wantDelta {
		t.Errorf("Expected delta to be nor more than %f, got %f", wantDelta, gotDelta)
	}
}

func getCreateGameData() game.CreateGame {
	return game.CreateGame{
		Name:        td.String(),
		Developer:   td.String(),
		Publisher:   td.String(),
		ReleaseDate: td.Date().Format("2006-01-02"),
		Genre:       []string{td.String(), td.String()},
		LogoURL:     td.String(),
	}
}

func compareCreateGameAndGame(t *testing.T, want game.CreateGame, got game.Game) {
	if want.Name != got.Name {
		t.Errorf("Expected to get game with name %s, got %s", want.Name, got.Name)
	}
	if want.Developer != got.Developer {
		t.Errorf("Expected to get game with developer %s, got %s", want.Developer, got.Developer)
	}
	if want.Publisher != got.Publisher {
		t.Errorf("Expected to get game with publisher %s, got %s", want.Publisher, got.Publisher)
	}
	if want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to get game with release date %s, got %s", want.ReleaseDate, got.ReleaseDate)
	}
	if len(want.Genre) != len(got.Genre) {
		t.Errorf("Expected to get game with %d genres, got %d", len(want.Genre), len(got.Genre))
	}
	if want.Genre[0] != got.Genre[0] {
		t.Errorf("Expected to get game with genre %s, got %s", want.Genre[0], got.Genre[0])
	}
	if want.LogoURL != got.LogoURL.String {
		t.Errorf("Expected to get game with logo url %s, got %s", want.LogoURL, got.LogoURL.String)
	}
}

func compareUpdateGameAndGame(t *testing.T, want game.UpdateGame, got game.Game) {
	if want.Name != got.Name {
		t.Errorf("Expected to get game with name %s, got %s", want.Name, got.Name)
	}
	if want.Developer != got.Developer {
		t.Errorf("Expected to get game with developer %s, got %s", want.Developer, got.Developer)
	}
	if want.Publisher != got.Publisher {
		t.Errorf("Expected to get game with publisher %s, got %s", want.Publisher, got.Publisher)
	}
	if want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to get game with release date %s, got %s", want.ReleaseDate, got.ReleaseDate)
	}
	if len(want.Genre) != len(got.Genre) {
		t.Errorf("Expected to get game with %d genres, got %d", len(want.Genre), len(got.Genre))
	}
	if want.Genre[0] != got.Genre[0] {
		t.Errorf("Expected to get game with genre %s, got %s", want.Genre[0], got.Genre[0])
	}
	if want.LogoURL != got.LogoURL.String {
		t.Errorf("Expected to get game with logo url %s, got %s", want.LogoURL, got.LogoURL.String)
	}
}
