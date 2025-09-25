package model

// ModerationItem represents a moderation entity for API response
type ModerationItem struct {
	ID           int32  `json:"id"`
	ResultStatus string `json:"resultStatus"`
	Details      string `json:"details"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
}
