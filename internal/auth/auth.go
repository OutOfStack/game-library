package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/OutOfStack/game-library/internal/client/authapi"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	// SigningAlgorithm specifies the default JWT signing algorithm (RS256)
	SigningAlgorithm = "RS256"
)

var tracer = otel.Tracer("auth")

var (
	tokenVerificationFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auth_token_verification_failures_total",
		Help: "Total number of token verification failures",
	})
)

// APIClient auth client api interface
type APIClient interface {
	VerifyToken(ctx context.Context, token string) (authapi.VerifyTokenResp, error)
}

// Client represents auth client
type Client struct {
	log           *zap.Logger
	parser        *jwt.Parser
	authAPIClient APIClient
}

// New constructs Auth instance
func New(log *zap.Logger, authAPIClient APIClient) (*Client, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{SigningAlgorithm}))

	return &Client{
		log:           log,
		parser:        parser,
		authAPIClient: authAPIClient,
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

// Verify calls Verify API method and returns error if token is invalid.
// If verify API is unavailable ErrVerifyAPIUnavailable is returned
func (c *Client) Verify(ctx context.Context, tokenStr string) error {
	ctx, span := tracer.Start(ctx, "verify")
	defer span.End()

	result, err := c.authAPIClient.VerifyToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	if !result.Valid {
		tokenVerificationFailures.Inc()
		return errors.New("invalid token")
	}

	return nil
}

// ParseToken returns token as a set of claims.
// No verification is done here, use Verify for that
func (c *Client) ParseToken(tokenStr string) (*Claims, error) {
	var claims Claims
	_, _, err := c.parser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	return &claims, nil
}
