package model_test

import (
	"encoding/json"
	"testing"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

func TestNewCreateModeration(t *testing.T) {
	gameID := td.Int31()
	data := model.ModerationData{Name: td.String(), Screenshots: []string{td.String()}}
	got := model.NewCreateModeration(gameID, data)
	require.Equal(t, gameID, got.GameID)
	require.Equal(t, data, got.GameData)
	require.Equal(t, model.ModerationStatusPending, got.Status)
}

func TestModerationData_ValueAndScan(t *testing.T) {
	tests := []struct {
		name string
		data model.ModerationData
	}{
		{name: "zero value"},
		{name: "populated", data: model.ModerationData{
			Name: "Game \"One\" Étoile", Developers: []string{"Studio"}, Publisher: "Publisher",
			ReleaseDate: "2026-09-05", Genres: []string{"Adventure"}, LogoURL: "https://example.com/cover.jpg",
			Summary: "First line\nSecond line", Screenshots: []string{"https://example.com/shot.jpg"},
			Websites: []string{"https://example.com"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.data.Value()
			require.NoError(t, err)
			text, ok := value.(string)
			require.True(t, ok)
			var decoded model.ModerationData
			require.NoError(t, json.Unmarshal([]byte(text), &decoded))
			require.Equal(t, tt.data, decoded)
			for _, input := range []any{text, []byte(text)} {
				var got model.ModerationData
				require.NoError(t, got.Scan(input))
				require.Equal(t, tt.data, got)
			}
		})
	}
}

func TestModerationData_Scan_InvalidInput(t *testing.T) {
	for _, input := range []any{42, "{", []byte("{"), `{"name":42}`} {
		var got model.ModerationData
		require.Error(t, got.Scan(input))
	}
}

func TestModerationData_Scan_NilPreservesValue(t *testing.T) {
	want := model.ModerationData{Name: td.String(), Genres: []string{td.String()}}
	got := want
	require.NoError(t, got.Scan(nil))
	require.Equal(t, want, got)
}
