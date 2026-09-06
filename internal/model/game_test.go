package model_test

import (
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestGetGameSlug(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "", want: ""},
		{name: "The Legend Of Zelda", want: "the-legend-of-zelda"},
		{name: "ÉTOILE CAFÉ", want: "étoile-café"},
		{name: "Game\xff Name", want: "game-name"},
		{name: "Game: Part-II!", want: "game:-part-ii!"},
		{name: " A  B ", want: "-a--b-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, model.GetGameSlug(tt.name))
		})
	}
}

func TestUpdateGame_MapToUpdateGameData(t *testing.T) {
	original := model.Game{
		Name: "Original Game", DevelopersIDs: []int32{1}, PublishersIDs: []int32{2},
		ReleaseDate: types.Date{Year: 2024, Month: time.February, Day: 29},
		GenresIDs:   []int32{3}, LogoURL: "original.jpg", Summary: "Original summary", Slug: "original-game",
		PlatformsIDs: []int32{4}, Screenshots: []string{"original-shot.jpg"}, Websites: []string{"original.example"},
		ModerationStatus: model.ModerationStatusReady,
	}
	baseline := model.UpdateGameData{
		Name: "Original Game", DevelopersIDs: []int32{5}, PublishersIDs: []int32{2},
		ReleaseDate: "2024-02-29", GenresIDs: []int32{3}, LogoURL: "original.jpg",
		Summary: "Original summary", Slug: "original-game", PlatformsIDs: []int32{4},
		Screenshots: []string{"original-shot.jpg"}, Websites: []string{"original.example"},
		ModerationStatus: model.ModerationStatusPending,
	}
	name, date, logo, summary := "New Game", "2026-09-05", "new.jpg", "New summary"
	genres, platforms := []int32{6, 7}, []int32{8, 9}
	screenshots, websites := []string{"new-shot.jpg"}, []string{"new.example"}
	empty := ""
	emptyIDs, emptyStrings := []int32{}, []string{}
	var nilIDs []int32
	var nilStrings []string
	tests := []struct {
		name    string
		request model.UpdateGame
		adjust  func(*model.UpdateGameData)
	}{
		{name: "omitted fields preserve original"},
		{
			name: "name changes slug only", request: model.UpdateGame{Name: &name},
			adjust: func(want *model.UpdateGameData) { want.Name, want.Slug = "New Game", "new-game" },
		},
		{
			name: "all fields replaced",
			request: model.UpdateGame{
				Name: &name, ReleaseDate: &date, GenresIDs: &genres, LogoURL: &logo, Summary: &summary,
				PlatformsIDs: &platforms, Screenshots: &screenshots, Websites: &websites,
			},
			adjust: func(want *model.UpdateGameData) {
				want.Name, want.Slug, want.ReleaseDate = "New Game", "new-game", "2026-09-05"
				want.GenresIDs, want.PlatformsIDs = []int32{6, 7}, []int32{8, 9}
				want.LogoURL, want.Summary = "new.jpg", "New summary"
				want.Screenshots, want.Websites = []string{"new-shot.jpg"}, []string{"new.example"}
			},
		},
		{
			name: "empty logo preserves original", request: model.UpdateGame{LogoURL: &empty},
		},
		{
			name: "explicit empty values clear fields",
			request: model.UpdateGame{
				Name: &empty, ReleaseDate: &empty, Summary: &empty, GenresIDs: &emptyIDs,
				PlatformsIDs: &emptyIDs, Screenshots: &emptyStrings, Websites: &emptyStrings,
			},
			adjust: func(want *model.UpdateGameData) {
				want.Name, want.Slug, want.ReleaseDate, want.Summary = "", "", "", ""
				want.GenresIDs, want.PlatformsIDs = []int32{}, []int32{}
				want.Screenshots, want.Websites = []string{}, []string{}
			},
		},
		{
			name:    "explicit nil slices clear fields",
			request: model.UpdateGame{GenresIDs: &nilIDs, PlatformsIDs: &nilIDs, Screenshots: &nilStrings, Websites: &nilStrings},
			adjust: func(want *model.UpdateGameData) {
				want.GenresIDs, want.PlatformsIDs, want.Screenshots, want.Websites = nil, nil, nil, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := baseline
			if tt.adjust != nil {
				tt.adjust(&want)
			}
			before := original
			got := tt.request.MapToUpdateGameData(original, []int32{5})
			require.Equal(t, want, got)
			require.Equal(t, before, original)
		})
	}
}

func TestCreateGame_MapToCreateGameData(t *testing.T) {
	input := model.CreateGame{
		Name: "New Game", Developer: "Studio", Publisher: "Publisher", ReleaseDate: "2026-09-05",
		GenresIDs: []int32{3}, LogoURL: "cover.jpg", Summary: "Summary", Slug: "custom-slug",
		PlatformsIDs: []int32{4}, Screenshots: []string{"shot.jpg"}, Websites: []string{"game.example"},
	}
	want := model.CreateGameData{
		Name: "New Game", DevelopersIDs: []int32{1}, PublishersIDs: []int32{2}, ReleaseDate: "2026-09-05",
		GenresIDs: []int32{3}, LogoURL: "cover.jpg", Summary: "Summary", Slug: "custom-slug",
		PlatformsIDs: []int32{4}, Screenshots: []string{"shot.jpg"}, Websites: []string{"game.example"},
		ModerationStatus: model.ModerationStatusPending,
	}
	require.Equal(t, want, input.MapToCreateGameData(2, 1))
}
