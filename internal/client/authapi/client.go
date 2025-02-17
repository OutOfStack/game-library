package authapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("")

// ErrVerifyAPIUnavailable - error representing unavailability of verify api
var ErrVerifyAPIUnavailable = errors.New("verify API is unavailable")

// Client represents dependencies for auth client
type Client struct {
	log    *zap.Logger
	apiURL string
	client *http.Client
}

// New constructs Client instance
func New(log *zap.Logger, apiURL string) (*Client, error) {
	client := &http.Client{
		Transport: observability.NewMonitoredTransport(otelhttp.NewTransport(http.DefaultTransport), "game-library-auth"),
	}

	return &Client{
		log:    log,
		apiURL: apiURL,
		client: client,
	}, nil
}

// VerifyToken returns result of token verification
func (c *Client) VerifyToken(ctx context.Context, token string) (VerifyTokenResp, error) {
	ctx, span := tracer.Start(ctx, "auth.verifyToken")
	defer span.End()

	data := VerifyToken{
		Token: token,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return VerifyTokenResp{}, fmt.Errorf("marshal verify token body: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewBuffer(body))
	if err != nil {
		return VerifyTokenResp{}, fmt.Errorf("create verify request: %v", err)
	}
	request.Header["Content-Type"] = []string{"application/json"}

	resp, err := c.client.Do(request)
	if err != nil {
		c.log.Error("call verify api", zap.String("url", c.apiURL), zap.Error(err))
		return VerifyTokenResp{}, ErrVerifyAPIUnavailable
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.String("url", c.apiURL), zap.Error(err))
		}
	}()

	var respBody VerifyTokenResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return VerifyTokenResp{}, fmt.Errorf("invalid response: %v", err)
	}

	return respBody, nil
}
