package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

// Context keys for authentication and authorization
const (
	CtxTokenKey  string = "auth_tkn"
	CtxClaimsKey string = "auth_clms"
)

// User roles
const (
	RoleModerator      = "moderator"
	RolePublisher      = "publisher"
	RoleRegisteredUser = "user"
)

var tracer = otel.Tracer("")

// ErrVerifyAPIUnavailable - error representing unavailability of verify api
var ErrVerifyAPIUnavailable = errors.New("verify API is unavailable")

// Claims represents jwt claims
type Claims struct {
	jwt.RegisteredClaims
	UserRole string `json:"user_role,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
}

// VerifyToken is a request type for JWT verification
type VerifyToken struct {
	Token string `json:"token" validate:"jwt"`
}

// VerifyTokenResp is a response type for JWT verification
type VerifyTokenResp struct {
	Valid bool `json:"valid"`
}

// Auth represents dependencies for auth methods
type Auth struct {
	log          *log.Logger
	parser       *jwt.Parser
	verifyAPIURL string
}

// New constructs Auth instance
func New(log *log.Logger, algorithm string, verifyAPIURL string) (*Auth, error) {
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm)
	}

	parser := &jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := Auth{
		log:          log,
		parser:       parser,
		verifyAPIURL: verifyAPIURL,
	}

	return &a, nil
}

// ExtractToken extracts Bearer token from Authorization hrader
func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}

	return ""
}

// Verify calls Verify API and returns nil if token is valid and error otherwise
// If verify API is unavailable ErrVerifyAPIUnavailable is returned
func (a *Auth) Verify(ctx context.Context, tokenStr string) error {
	ctx, span := tracer.Start(ctx, "auth.verify")
	defer span.End()

	data := VerifyToken{
		Token: tokenStr,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshalling verify token body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", a.verifyAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating verify request")
	}
	request.Header["Content-Type"] = []string{"application/json"}

	resp, err := otelhttp.DefaultClient.Do(request)
	if err != nil {
		log.Printf("error calling verify api at %s: %v\n", a.verifyAPIURL, err)
		return ErrVerifyAPIUnavailable
	}
	defer resp.Body.Close()

	var respBody VerifyTokenResp
	json.NewDecoder(resp.Body).Decode(&respBody)

	if respBody.Valid {
		return nil
	}
	return fmt.Errorf("invalid token")
}

// ParseToken returns token as a set of claims
func (a *Auth) ParseToken(tokenStr string) (*Claims, error) {
	var claims Claims
	_, _, err := a.parser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		return &Claims{}, fmt.Errorf("parsing token: %w", err)
	}
	return &claims, nil
}
