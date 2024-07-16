package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// Context keys for authentication and authorization
const (
	CtxTokenKey  string = "auth_token"
	CtxClaimsKey string = "auth_claims"
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

// UserID return user id from claims
func (c *Claims) UserID() string {
	return c.Subject
}

// VerifyToken is a request type for JWT verification
type VerifyToken struct {
	Token string `json:"token" validate:"jwt"`
}

// VerifyTokenResp is a response type for JWT verification
type VerifyTokenResp struct {
	Valid bool `json:"valid"`
}

// Client represents auth client
type Client struct {
	log          *zap.Logger
	parser       *jwt.Parser
	verifyAPIURL string
}

// New constructs Auth instance
func New(log *zap.Logger, algorithm string, verifyAPIURL string) (*Client, error) {
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm)
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{algorithm}))

	return &Client{
		log:          log,
		parser:       parser,
		verifyAPIURL: verifyAPIURL,
	}, nil
}

// ExtractToken extracts Bearer token from Authorization header
func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}

	return ""
}

// Verify calls Verify API and returns error if token is invalid.
// If verify API is unavailable ErrVerifyAPIUnavailable is returned
func (c *Client) Verify(ctx context.Context, tokenStr string) error {
	ctx, span := tracer.Start(ctx, "auth.verify")
	defer span.End()

	data := VerifyToken{
		Token: tokenStr,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal verify token body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", c.verifyAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create verify request: %v", err)
	}
	request.Header["Content-Type"] = []string{"application/json"}

	resp, err := otelhttp.DefaultClient.Do(request)
	if err != nil {
		c.log.Error("call verify api", zap.String("url", c.verifyAPIURL), zap.Error(err))
		return ErrVerifyAPIUnavailable
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.String("url", c.verifyAPIURL), zap.Error(err))
		}
	}()

	var respBody VerifyTokenResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return fmt.Errorf("invalid response: %v", err)
	}

	if !respBody.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

// ParseToken returns token as a set of claims
func (c *Client) ParseToken(tokenStr string) (*Claims, error) {
	var claims Claims
	_, _, err := c.parser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		return &Claims{}, fmt.Errorf("parsing token: %w", err)
	}
	return &claims, nil
}
