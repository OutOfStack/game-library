package openaiapi

// ModerationResponse represents OpenAI moderation API response
type ModerationResponse struct {
	ID      string             `json:"id"`
	Results []ModerationResult `json:"results"`
}

// ModerationResult represents moderation result for single input
type ModerationResult struct {
	Flagged    bool     `json:"flagged"`
	Categories []string `json:"categories"`
}

// VisionAnalysisResult represents the result from vision model image analysis
type VisionAnalysisResult struct {
	Approved          bool   `json:"approved"`
	Reason            string `json:"reason"`
	GamingAppropriate bool   `json:"gaming_appropriate"`
	ContentRelevant   bool   `json:"content_relevant"`
}
