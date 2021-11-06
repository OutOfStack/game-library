package game_test

import (
	"context"
	"errors"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
)

// TestGetInfos_NoData_Empty tests case when there are no games in table
func TestGetInfos_NoData_Empty(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	games, err := game.GetInfos(context.Background(), db, 20, 0)
	if err != nil {
		t.Fatalf("error getting games: %v", err)
	}

	want := 0
	got := len(games)
	if want != got {
		t.Fatalf("Expected to retrieve %d games, got %d", want, got)
	}
}

// TestGetInfos_DataExists_Equals tests case when we add game and then get games from database and compare them
func TestGetInfos_DataExists_Equals(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	g := game.CreateGameReq{
		Name:        "Test game",
		Developer:   "Test developer",
		Publisher:   "Test publisher",
		ReleaseDate: "2021-11-03",
		Price:       100,
		Genre:       []string{"rpg"},
		LogoUrl:     "http://images/123",
	}

	_, err := game.Create(context.Background(), db, g)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	games, err := game.GetInfos(context.Background(), db, 20, 0)
	if err != nil {
		t.Fatalf("error getting games: %v", err)
	}

	if len(games) == 0 {
		t.Fatal("Expected to retrieve 1 games, got 0")
	}

	want := g
	got := games[0]
	if want.Name != got.Name {
		t.Errorf("Expected to retrieve game with name %s, got %s", want.Name, got.Name)
	}
	if want.Developer != got.Developer {
		t.Errorf("Expected to retrieve game with developer %s, got %s", want.Developer, got.Developer)
	}
	if want.Publisher != got.Publisher {
		t.Errorf("Expected to retrieve game with publisher %s, got %s", want.Publisher, got.Publisher)
	}
	if want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to retrieve game with release date %s, got %s", want.ReleaseDate, got.ReleaseDate)
	}
	if want.Price != got.Price {
		t.Errorf("Expected to retrieve game with price %f, got %f", want.Price, got.Price)
	}
	if want.Genre[0] != got.Genre[0] {
		t.Errorf("Expected to retrieve game with genre %s, got %s", want.Genre[0], got.Genre[0])
	}
	if want.LogoUrl != got.LogoUrl.String {
		t.Errorf("Expected to retrieve game with logo url %s, got %s", want.LogoUrl, got.LogoUrl.String)
	}
}

// TestRetrieveInfo_NoData_Empty tests case when there is no game with provided id
func TestRetrieveInfo_NoData_Empty(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	var id int64 = 1234
	g, err := game.RetrieveInfo(context.Background(), db, id)
	if err != nil {
		if errors.Is(err, game.ErrNotFound{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = game.ErrNotFound{
		Entity: "game",
		ID:     id,
	}
	gotErr := err
	if g != nil || !errors.Is(gotErr, wantErr) {
		t.Fatalf("Expected to receive error %v, got %v", wantErr, gotErr)
	}
}

// TestRetrieveInfo_DataExists_Equals tests case when we add game and then get this game by id from database and compare them
func TestRetrieveInfo_DataExists_Equals(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	g := game.CreateGameReq{
		Name:        "Test game",
		Developer:   "Test developer",
		Publisher:   "Test publisher",
		ReleaseDate: "2021-11-03",
		Price:       100,
		Genre:       []string{"rpg"},
		LogoUrl:     "http://images/123",
	}

	id, err := game.Create(context.Background(), db, g)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	gg, err := game.RetrieveInfo(context.Background(), db, id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	want := g
	got := gg
	if want.Name != got.Name {
		t.Errorf("Expected to retrieve game with name %s, got %s", want.Name, got.Name)
	}
	if want.Developer != got.Developer {
		t.Errorf("Expected to retrieve game with developer %s, got %s", want.Developer, got.Developer)
	}
	if want.Publisher != got.Publisher {
		t.Errorf("Expected to retrieve game with publisher %s, got %s", want.Publisher, got.Publisher)
	}
	if want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to retrieve game with release date %s, got %s", want.ReleaseDate, got.ReleaseDate)
	}
	if want.Price != got.Price {
		t.Errorf("Expected to retrieve game with price %f, got %f", want.Price, got.Price)
	}
	if want.Genre[0] != got.Genre[0] {
		t.Errorf("Expected to retrieve game with genre %s, got %s", want.Genre[0], got.Genre[0])
	}
	if want.LogoUrl != got.LogoUrl.String {
		t.Errorf("Expected to retrieve game with logo url %s, got %s", want.LogoUrl, got.LogoUrl.String)
	}
}
