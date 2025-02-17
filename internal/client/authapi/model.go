package authapi

// VerifyToken is a request type for JWT verification
type VerifyToken struct {
	Token string `json:"token" validate:"jwt"`
}

// VerifyTokenResp is a response type for JWT verification
type VerifyTokenResp struct {
	Valid bool `json:"valid"`
}
