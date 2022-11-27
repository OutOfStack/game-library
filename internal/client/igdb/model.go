package igdb

// TokenResp - access token response
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TopRatedGamesResp - top rated games response
type TopRatedGamesResp struct {
	ID          uint32  `json:"id"`
	Name        string  `json:"name"`
	Rating      float32 `json:"rating"`
	RatingCount uint32  `json:"rating_count"`
}
