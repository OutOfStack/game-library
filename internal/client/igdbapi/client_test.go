package igdbapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetTopRatedGames_RequestsAndDecodesRichFields(t *testing.T) {
	client, requests := newTestClient(t, `[{"id":1,"videos":[{"video_id":"trailer"}],"game_modes":[1,3],"similar_games":[2],"franchise":{"games":[3]},"franchises":[{"games":[4]}]}]`)

	games, err := client.GetTopRatedGames(t.Context(), []int64{6}, time.Unix(100, 0), 1, 2, 3)
	require.NoError(t, err)
	require.Len(t, games, 1)
	require.Equal(t, "trailer", games[0].Videos[0].VideoID)
	require.Equal(t, []int32{1, 3}, games[0].GameModes)
	require.Equal(t, []int64{2}, games[0].SimilarGames)
	require.Equal(t, []int64{3}, games[0].Franchise.Games)
	require.Equal(t, []int64{4}, games[0].Franchises[0].Games)
	require.Contains(t, (*requests)[0], "videos.video_id")
	require.Contains(t, (*requests)[0], "game_modes")
	require.Contains(t, (*requests)[0], "similar_games")
	require.Contains(t, (*requests)[0], "franchise.games")
	require.Contains(t, (*requests)[0], "franchises.games")
}

func TestGetGameInfoForUpdate_RequestsAndDecodesRichFields(t *testing.T) {
	client, requests := newTestClient(t, `[{"id":1,"videos":[{"video_id":"trailer"}],"game_modes":[2],"similar_games":[5],"franchise":{"games":[6]},"franchises":[{"games":[7]}]}]`)

	game, err := client.GetGameInfoForUpdate(t.Context(), 1)
	require.NoError(t, err)
	require.Equal(t, "trailer", game.Videos[0].VideoID)
	require.Equal(t, []int32{2}, game.GameModes)
	require.Equal(t, []int64{5}, game.SimilarGames)
	require.Equal(t, []int64{6}, game.Franchise.Games)
	require.Equal(t, []int64{7}, game.Franchises[0].Games)
	require.Contains(t, (*requests)[0], "videos.video_id")
	require.Contains(t, (*requests)[0], "game_modes")
	require.Contains(t, (*requests)[0], "similar_games")
	require.Contains(t, (*requests)[0], "franchise.games")
	require.Contains(t, (*requests)[0], "franchises.games")
}

func newTestClient(t *testing.T, response string) (*Client, *[]string) {
	t.Helper()

	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		requests = append(requests, string(body))
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(response))
		if err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)

	return &Client{
		log:        zap.NewNop(),
		clientID:   "client",
		apiURL:     server.URL,
		token:      &tokenInfo{token: "token", expiresAt: time.Now().Add(time.Hour)},
		httpClient: server.Client(),
	}, &requests
}
