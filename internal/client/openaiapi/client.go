package openaiapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	// ai-related processing is longer that other generic apis
	defaultTimeout = 30 * time.Second

	gameSummaryMaxLen = 2000

	moderationEndpoint = "moderations"
	visionEndpoint     = "/chat/completions"

	maxVisionTokens = 500
)

var tracer = otel.Tracer("openai")

// Client represents OpenAI API client
type Client struct {
	log             *zap.Logger
	httpClient      *http.Client
	apiKey          string
	apiURL          string
	moderationModel string
	visionModel     string
}

// New creates new OpenAI client
func New(log *zap.Logger, cfg appconf.OpenAI) *Client {
	httpClient := &http.Client{
		Transport: observability.NewMonitoredTransport(otelhttp.NewTransport(http.DefaultTransport), "openai"),
		Timeout:   defaultTimeout,
	}

	return &Client{
		httpClient:      httpClient,
		apiKey:          cfg.APIKey,
		apiURL:          cfg.APIURL,
		moderationModel: cfg.ModerationModel,
		visionModel:     cfg.VisionModel,
		log:             log,
	}
}

// ModerateText performs basic text and image moderation using OpenAI moderation API
func (c *Client) ModerateText(ctx context.Context, gameData model.ModerationData) (*ModerationResponse, error) {
	ctx, span := tracer.Start(ctx, "ModerateText")
	defer span.End()

	req := c.getModerationRequest(gameData)

	span.SetAttributes(
		attribute.String("game.name", gameData.Name),
		attribute.String("openai.model", req.Model),
		attribute.Int("openai.input_count", len(req.Input)))

	reqURL, err := url.JoinPath(c.apiURL, moderationEndpoint)
	if err != nil {
		return nil, fmt.Errorf("join url path: %v", err)
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.String("url", reqURL), zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var moderationResp ModerationResponse
	if err = json.NewDecoder(resp.Body).Decode(&moderationResp); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}

	return &moderationResp, nil
}

// AnalyzeGameImages analyzes images for gaming-specific content appropriateness using vision model
func (c *Client) AnalyzeGameImages(ctx context.Context, gameData model.ModerationData) (*VisionAnalysisResult, error) {
	ctx, span := tracer.Start(ctx, "AnalyzeGameImages")
	defer span.End()

	req := c.getVisionRequest(gameData)

	span.SetAttributes(
		attribute.String("game.name", gameData.Name),
		attribute.String("openai.model", req.Model),
		attribute.Int("openai.max_tokens", req.MaxTokens))

	reqURL, err := url.JoinPath(c.apiURL, visionEndpoint)
	if err != nil {
		return nil, fmt.Errorf("join url path: %v", err)
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.String("url", reqURL), zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var visionResp VisionResponse
	if err = json.NewDecoder(resp.Body).Decode(&visionResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return parseVisionResponse(&visionResp)
}

func (c *Client) getModerationRequest(gameData model.ModerationData) ModerationRequest {
	// text inputs
	inputs := []ModerationInputItem{
		{Type: textType, Text: gameData.Name},
		{Type: textType, Text: truncateText(gameData.Summary, gameSummaryMaxLen)},
		{Type: textType, Text: strings.Join(gameData.Developers, ", ")},
		{Type: textType, Text: gameData.Publisher},
		{Type: textType, Text: strings.Join(gameData.Websites, ", ")},
	}

	// image inputs
	if gameData.LogoURL != "" {
		inputs = append(inputs, ModerationInputItem{
			Type:     imageURLType,
			ImageURL: gameData.LogoURL,
		})
	}

	for _, scr := range gameData.Screenshots {
		inputs = append(inputs, ModerationInputItem{
			Type:     imageURLType,
			ImageURL: scr,
		})
	}

	return ModerationRequest{
		Model: c.moderationModel,
		Input: inputs,
	}
}

func (c *Client) getVisionRequest(gameData model.ModerationData) VisionRequest {
	// prompt
	prompt := buildGamingModerationPrompt(gameData)

	// images
	images := make([]string, 0, 1+len(gameData.Screenshots))
	if gameData.LogoURL != "" {
		images = append(images, gameData.LogoURL)
	}

	images = append(images, gameData.Screenshots...)

	// vision request
	content := []VisionContent{
		{Type: textType, Text: prompt},
	}

	for _, imageURL := range images {
		content = append(content, VisionContent{
			Type:     imageType,
			ImageURL: &VisionImageURL{URL: imageURL},
		})
	}

	return VisionRequest{
		Model:     c.visionModel,
		MaxTokens: maxVisionTokens,
		Messages: []VisionMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}
}

// buildGamingModerationPrompt creates a gaming-specific moderation prompt
func buildGamingModerationPrompt(gameData model.ModerationData) string {
	return fmt.Sprintf(`You are moderating content for a video game library platform. Analyze the provided images and context for appropriateness.

Game Context:
- Name: %s
- Genre: %s
- Summary: %s
- Publisher: %s

Gaming Content Guidelines:
1. ALLOWED: Typical video game violence (shooting, fighting, fantasy combat) - this is NORMAL for games
2. ALLOWED: Video game weapons, explosions, action scenes - expected in action games
3. ALLOWED: Stylized/cartoon violence, sci-fi themes, fantasy elements
4. FLAGGED: Extremely graphic realistic violence with excessive blood/gore
5. FLAGGED: Real-world hate symbols, explicit sexual content, illegal activities
6. FLAGGED: Images completely unrelated to gaming (random photos, spam content)
7. FLAGGED: Personal information, contact details, or promotional spam

Be GAMING-FRIENDLY - most action game content should be approved unless extremely inappropriate.

Respond ONLY with JSON:
{
  "approved": true/false,
  "reason": "brief explanation",
  "gaming_appropriate": true/false,
  "content_relevant": true/false
}`, gameData.Name, strings.Join(gameData.Genres, ", "), gameData.Summary, gameData.Publisher)
}

// parseVisionResponse extracts moderation result from vision API response
func parseVisionResponse(resp *VisionResponse) (*VisionAnalysisResult, error) {
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices in response")
	}

	content := resp.Choices[0].Message.Content
	if len(content) == 0 {
		return nil, errors.New("no content in response")
	}

	var result VisionAnalysisResult
	text := content[0].Text
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parse JSON response: %w", err)
	}

	return &result, nil
}

func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
