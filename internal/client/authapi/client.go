package authapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	defaultTimeout = 10 * time.Second
)

var tracer = otel.Tracer("authapi")

// ErrVerifyAPIUnavailable - error representing unavailability of verify api
var ErrVerifyAPIUnavailable = errors.New("verify API is unavailable")

// Client represents dependencies for auth client
type Client struct {
	log                    *zap.Logger
	httpClient             *http.Client
	verifyTokenEndpointURL string
}

// New constructs Client instance
func New(log *zap.Logger, verifyTokenEndpointURL string) (*Client, error) {
	client := &http.Client{
		Transport: observability.NewMonitoredTransport(otelhttp.NewTransport(http.DefaultTransport), "game-library-auth"),
		Timeout:   defaultTimeout,
	}

	return &Client{
		log:                    log,
		verifyTokenEndpointURL: verifyTokenEndpointURL,
		httpClient:             client,
	}, nil
}

// VerifyToken returns result of token verification
func (c *Client) VerifyToken(ctx context.Context, token string) (VerifyTokenResp, error) {
	ctx, span := tracer.Start(ctx, "verifyToken")
	defer span.End()

	data := VerifyToken{
		Token: token,
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		return VerifyTokenResp{}, fmt.Errorf("marshal verify token body: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.verifyTokenEndpointURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return VerifyTokenResp{}, fmt.Errorf("create verify request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.log.Error("call verify api", zap.String("url", c.verifyTokenEndpointURL), zap.Error(err))
		return VerifyTokenResp{}, ErrVerifyAPIUnavailable
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.String("url", c.verifyTokenEndpointURL), zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return VerifyTokenResp{}, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var respBody VerifyTokenResp
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return VerifyTokenResp{}, fmt.Errorf("invalid response: %v", err)
	}

	return respBody, nil
}
