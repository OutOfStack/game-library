package openaiapi

const (
	textType     = "text"
	imageType    = "image"
	imageURLType = "image_url"
)

// ModerationRequest represents OpenAI moderation API request
type ModerationRequest struct {
	Model string                `json:"model"`
	Input []ModerationInputItem `json:"input"`
}

// ModerationInputItem represents input item for moderation
type ModerationInputItem struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// ModerationResponse represents OpenAI moderation API response
type ModerationResponse struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

// ModerationResult represents moderation result for single input
type ModerationResult struct {
	Flagged        bool               `json:"flagged"`
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
}

// VisionRequest represents request to vision model for complex analysis
type VisionRequest struct {
	Model          string            `json:"model"`
	Messages       []VisionMessage   `json:"messages"`
	MaxTokens      int               `json:"max_tokens"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

// VisionMessage represents message in vision request
type VisionMessage struct {
	Role    string          `json:"role"`
	Content []VisionContent `json:"content"`
}

// VisionContent represents content in vision message
type VisionContent struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *VisionImageURL `json:"image_url,omitempty"`
}

// VisionImageURL represents image URL in vision content
type VisionImageURL struct {
	URL string `json:"url"`
}

// VisionResponse represents vision model response
type VisionResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []VisionChoice `json:"choices"`
}

// VisionChoice represents choice in vision response
type VisionChoice struct {
	Index   int           `json:"index"`
	Message VisionMessage `json:"message"`
}

// GameModerationResult represents the final moderation result for a game
type GameModerationResult struct {
	Approved       bool     `json:"approved"`
	Reason         string   `json:"reason"`
	Details        string   `json:"details"`
	ViolationTypes []string `json:"violation_types,omitempty"`
}

// VisionAnalysisResult represents the result from vision model image analysis
type VisionAnalysisResult struct {
	Approved          bool   `json:"approved"`
	Reason            string `json:"reason"`
	GamingAppropriate bool   `json:"gaming_appropriate"`
	ContentRelevant   bool   `json:"content_relevant"`
}
