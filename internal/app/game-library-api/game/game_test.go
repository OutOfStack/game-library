package game_test

import (
	"context"
	"errors"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
)

var newGame game.CreateGameReq = game.CreateGameReq{
	Name:        "Test game",
	Developer:   "Test developer",
	Publisher:   "Test publisher",
	ReleaseDate: "2021-11-03",
	Price:       100,
	Genre:       []string{"rpg"},
	LogoUrl:     "http://images/123",
}

// TestGetInfos_NotExist_ShouldReturnEmpty tests case when there is no data and we should get empty result
func TestGetInfos_NotExist_ShouldReturnEmpty(t *testing.T) {
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

// TestGetInfos_DataExists_ShouldBeEqual tests case when we add one game, then fetch first game and they should be equal
func TestGetInfos_DataExists_ShouldBeEqual(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	g := newGame

	_, err := game.Create(context.Background(), db, g)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	games, err := game.GetInfos(context.Background(), db, 20, 0)
	if err != nil {
		t.Fatalf("error getting games: %v", err)
	}

	if len(games) == 0 {
		t.Fatal("Expected to retrieve 1 game, got 0")
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

// TestRetrieveInfo_NotExist_ShouldReturnNotFoundError tests case when a game with provided id does not exist and we should get a Not Found Error
func TestRetrieveInfo_NotExist_ShouldReturnNotFoundError(t *testing.T) {
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
		if !errors.As(err, &game.ErrNotFound{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = game.ErrNotFound{
		Entity: "game",
		ID:     id,
	}
	gotErr := err
	if g != nil || !errors.Is(gotErr, wantErr) {
		t.Fatalf("Expected to receive empty entity and error [%v], got [%v]", wantErr, gotErr)
	}
}

// TestRetrieveInfo_DataExists_ShouldRetrieveEqual tests case when we add game, then fetch this game and they should be equal
func TestRetrieveInfo_DataExists_ShouldRetrieveEqual(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	g := newGame

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

// TestSearchInfos_DataExists_ShouldReturnEqual tests case when we add game, then search this game and they should be equal
func TestSearchInfos_DataExists_ShouldReturnEqual(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	ng := newGame

	_, err := game.Create(context.Background(), db, ng)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	matched, err := game.SearchInfos(context.Background(), db, ng.Name)
	if err != nil {
		t.Fatalf("error searching games: %v", err)
	}

	if len(matched) == 0 {
		t.Fatal("Expected to retrieve 1 matched game, got 0")
	}

	want := ng
	got := matched[0]
	if want.Name != got.Name {
		t.Errorf("Expected to find game with name %s, got %s", want.Name, got.Name)
	}
	if want.Developer != got.Developer {
		t.Errorf("Expected to find game with developer %s, got %s", want.Developer, got.Developer)
	}
	if want.Publisher != got.Publisher {
		t.Errorf("Expected to find game with publisher %s, got %s", want.Publisher, got.Publisher)
	}
	if want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to find game with release date %s, got %s", want.ReleaseDate, got.ReleaseDate)
	}
	if want.Price != got.Price {
		t.Errorf("Expected to find game with price %f, got %f", want.Price, got.Price)
	}
	if want.Genre[0] != got.Genre[0] {
		t.Errorf("Expected to find game with genre %s, got %s", want.Genre[0], got.Genre[0])
	}
	if want.LogoUrl != got.LogoUrl.String {
		t.Errorf("Expected to find game with logo url %s, got %s", want.LogoUrl, got.LogoUrl.String)
	}
}

// TestSearchInfos_DataExists_ShouldReturnMatched tests case when we add multiple games, then search games and we should get matched
func TestSearchInfos_DataExists_ShouldReturnMatched(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	ng1 := newGame
	ng1.Name = "test game name"
	ng2 := newGame
	ng2.Name = "tEsTGameName"
	ng3 := newGame
	ng3.Name = "tEssTGameName"
	ng4 := newGame
	ng4.Name = "a test game name"

	_, err := game.Create(context.Background(), db, ng1)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}
	_, err = game.Create(context.Background(), db, ng2)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}
	_, err = game.Create(context.Background(), db, ng3)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}
	_, err = game.Create(context.Background(), db, ng4)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	matched, err := game.SearchInfos(context.Background(), db, "test")
	if err != nil {
		t.Fatalf("error searching games: %v", err)
	}

	want := 2
	got := len(matched)
	if want != got {
		t.Fatalf("Expected to retrieve %d matched game, got %d", want, got)
	}
}

// TestUpdate_Valid_ShouldRetrieveEqual tests case when we update game, then fetch this game and they should be equal
func TestUpdate_Valid_ShouldRetrieveEqual(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	cr := newGame

	id, err := game.Create(context.Background(), db, cr)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	n, d, p, rd, pr, ge, lu := "New game", "New developer", "New publisher", "2021-11-12", float32(50), []string{"adventure"}, "https://images/999"
	up := game.UpdateGameReq{
		Name:        &n,
		Developer:   &d,
		Publisher:   &p,
		ReleaseDate: &rd,
		Price:       &pr,
		Genre:       &ge,
		LogoUrl:     &lu,
	}

	_, err = game.Update(context.Background(), db, id, up)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	gg, err := game.RetrieveInfo(context.Background(), db, id)
	if err != nil {
		t.Fatalf("error getting game: %v", err)
	}

	want := up
	got := gg
	if *want.Name != got.Name {
		t.Errorf("Expected to retrieve game with name %s, got %s", *want.Name, got.Name)
	}
	if *want.Developer != got.Developer {
		t.Errorf("Expected to retrieve game with developer %s, got %s", *want.Developer, got.Developer)
	}
	if *want.Publisher != got.Publisher {
		t.Errorf("Expected to retrieve game with publisher %s, got %s", *want.Publisher, got.Publisher)
	}
	if *want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to retrieve game with release date %s, got %s", *want.ReleaseDate, got.ReleaseDate)
	}
	if *want.Price != got.Price {
		t.Errorf("Expected to retrieve game with price %f, got %f", *want.Price, got.Price)
	}
	if (*want.Genre)[0] != got.Genre[0] {
		t.Errorf("Expected to retrieve game with genre %s, got %s", (*want.Genre)[0], got.Genre[0])
	}
	if *want.LogoUrl != got.LogoUrl.String {
		t.Errorf("Expected to retrieve game with logo url %s, got %s", *want.LogoUrl, got.LogoUrl.String)
	}
}

// TestUpdate_Valid_ShouldReturnEqual tests case when we update game and returned value should be equal
func TestUpdate_Valid_ShouldReturnEqual(t *testing.T) {
	if err := setup(db); err != nil {
		t.Fatalf("error on setup: %v", err)
	}

	defer func() {
		if err := teardown(db); err != nil {
			t.Fatalf("error on teardown: %v", err)
		}
	}()

	defer recovery(t)

	cr := newGame

	id, err := game.Create(context.Background(), db, cr)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	n, d, p, rd, pr, ge, lu := "New game", "New developer", "New publisher", "2021-11-12", float32(50), []string{"adventure"}, "https://images/999"
	up := game.UpdateGameReq{
		Name:        &n,
		Developer:   &d,
		Publisher:   &p,
		ReleaseDate: &rd,
		Price:       &pr,
		Genre:       &ge,
		LogoUrl:     &lu,
	}

	gg, err := game.Update(context.Background(), db, id, up)
	if err != nil {
		t.Fatalf("error creating game: %v", err)
	}

	want := up
	got := gg
	if *want.Name != got.Name {
		t.Errorf("Expected to return game with name %s, got %s", *want.Name, got.Name)
	}
	if *want.Developer != got.Developer {
		t.Errorf("Expected to return game with developer %s, got %s", *want.Developer, got.Developer)
	}
	if *want.Publisher != got.Publisher {
		t.Errorf("Expected to return game with publisher %s, got %s", *want.Publisher, got.Publisher)
	}
	if *want.ReleaseDate != got.ReleaseDate.String() {
		t.Errorf("Expected to return game with release date %s, got %s", *want.ReleaseDate, got.ReleaseDate)
	}
	if *want.Price != got.Price {
		t.Errorf("Expected to return game with price %f, got %f", *want.Price, got.Price)
	}
	if (*want.Genre)[0] != got.Genre[0] {
		t.Errorf("Expected to return game with genre %s, got %s", (*want.Genre)[0], got.Genre[0])
	}
	if *want.LogoUrl != got.LogoUrl.String {
		t.Errorf("Expected to return game with logo url %s, got %s", *want.LogoUrl, got.LogoUrl.String)
	}
}

// TestUpdate_NotExist_ShouldReturnNotFoundError tests case when we update a non existing game and we should get a Not Found Error
func TestUpdate_NotExist_ShouldReturnNotFoundError(t *testing.T) {
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
	up := game.UpdateGameReq{}
	g, err := game.Update(context.Background(), db, id, up)
	if err != nil {
		if !errors.As(err, &game.ErrNotFound{}) {
			t.Fatalf("error getting game: %v", err)
		}
	}

	var wantErr = game.ErrNotFound{
		Entity: "game",
		ID:     id,
	}
	gotErr := err
	if g != nil || !errors.Is(gotErr, wantErr) {
		t.Fatalf("Expected to receive empty entity and error [%v], got [%v]", wantErr, gotErr)
	}
}
