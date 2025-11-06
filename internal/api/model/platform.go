package model

// Platform represents platform response
type Platform struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}
