//go:build integration

package openaiapi_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/client/openaiapi"
	"github.com/OutOfStack/game-library/internal/model"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	defaultAPIURL          = "https://api.openai.com/v1"
	defaultModerationModel = "omni-moderation-latest"
	defaultVisionModel     = "gpt-5.4-nano"
	defaultTimeout         = 60 * time.Second

	// stable public test image (dog photo)
	testLogoURL = "https://picsum.photos/id/237/400/300"
)

// newTestClient creates OpenAI client configured from app.env in repo root
// (same source the app uses) with fallback to env variables,
// fails the test if OPENAI_API_KEY is not set
func newTestClient(t *testing.T) *openaiapi.Client {
	t.Helper()

	v := viper.New()
	envFile := filepath.Join("..", "..", "..", "app.env")
	if _, err := os.Stat(envFile); err == nil {
		v.SetConfigFile(envFile)
		v.SetConfigType("env")
		require.NoError(t, v.ReadInConfig())
	}
	v.AutomaticEnv()

	conf := appconf.OpenAI{
		APIKey:          v.GetString("OPENAI_API_KEY"),
		APIURL:          v.GetString("OPENAI_API_URL"),
		ModerationModel: v.GetString("OPENAI_MODERATION_MODEL"),
		VisionModel:     v.GetString("OPENAI_VISION_MODEL"),
		Timeout:         v.GetDuration("OPENAI_API_TIMEOUT"),
	}
	if conf.APIKey == "" {
		t.Fatal("OPENAI_API_KEY is not set (neither in app.env nor in env)")
	}
	if conf.APIURL == "" {
		conf.APIURL = defaultAPIURL
	}
	if conf.ModerationModel == "" {
		conf.ModerationModel = defaultModerationModel
	}
	if conf.VisionModel == "" {
		conf.VisionModel = defaultVisionModel
	}
	if conf.Timeout <= 0 {
		conf.Timeout = defaultTimeout
	}

	return openaiapi.New(zaptest.NewLogger(t), conf)
}

func TestModerateText_CleanContent(t *testing.T) {
	client := newTestClient(t)

	data := model.ModerationData{
		Name:       "Star Voyager",
		Summary:    "A relaxing space exploration game where players chart unknown star systems, trade with alien civilizations and upgrade their ship",
		Developers: []string{"Orbit Studio"},
		Publisher:  "Galactic Games",
		Websites:   []string{"https://example.com/star-voyager"},
		LogoURL:    testLogoURL,
	}

	resp, err := client.ModerateText(t.Context(), data)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	// moderation API returns a single aggregated result for a multimodal input array
	require.NotEmpty(t, resp.Results)
	for i, r := range resp.Results {
		assert.False(t, r.Flagged, "input %d unexpectedly flagged with categories: %v", i, r.Categories)
	}
}

func TestModerateText_FlaggedContent(t *testing.T) {
	client := newTestClient(t)

	tests := []struct {
		name             string
		summary          string
		expectedCategory string
	}{
		{
			name:             "violence",
			summary:          "I will hunt you down and kill you and your whole family, you worthless subhuman scum. You deserve to die.",
			expectedCategory: "violence",
		},
		{
			name:             "self-harm",
			summary:          "I want to kill myself, I am looking for the most painless way to end my own life tonight",
			expectedCategory: "self-harm",
		},
		{
			name:             "illicit",
			summary:          "A detailed step by step guide on how to cook methamphetamine at home and sell it to school students",
			expectedCategory: "illicit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := model.ModerationData{
				Name:       "Totally Legit Game",
				Summary:    tc.summary,
				Developers: []string{"Some Studio"},
				Publisher:  "Some Publisher",
			}

			resp, err := client.ModerateText(t.Context(), data)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp.Results)

			result := resp.Results[0]
			assert.True(t, result.Flagged, "expected content to be flagged")
			assert.Contains(t, result.Categories, tc.expectedCategory)
			assert.Contains(t, result.CategoryInputTypes[tc.expectedCategory], "text")
		})
	}
}

func TestAnalyzeGameImages(t *testing.T) {
	client := newTestClient(t)

	data := model.ModerationData{
		Name:      "Star Voyager",
		Genres:    []string{"Adventure", "Simulation"},
		Summary:   "A relaxing space exploration game where players chart unknown star systems and trade with alien civilizations",
		Publisher: "Galactic Games",
		LogoURL:   testLogoURL,
	}

	result, err := client.AnalyzeGameImages(t.Context(), data)
	require.NoError(t, err)
	require.NotNil(t, result)
	// vision verdict on a placeholder photo is model judgment and may vary,
	// so only verify the response is well-formed
	assert.NotEmpty(t, result.Reason)
	t.Logf("vision result: approved=%t, gaming_appropriate=%t, content_relevant=%t, reason=%q",
		result.Approved, result.GamingAppropriate, result.ContentRelevant, result.Reason)
}
