package igdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type token struct {
	token     string
	expiresAt time.Time
	mu        sync.RWMutex
}

func (t *token) get() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.expiresAt.After(time.Now().Add(5 * time.Minute)) {
		return t.token
	}
	return ""
}

func (t *token) set(token string, expiresIn int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.token = token
	t.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// accessToken returns access token
func (c *Client) accessToken(ctx context.Context) (string, error) {
	token := c.token.get()
	if token != "" {
		return token, nil
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.conf.TokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating verify request")
	}

	q := req.URL.Query()
	q.Add("grant_type", "client_credentials")
	q.Add("client_id", c.conf.ClientID)
	q.Add("client_secret", c.conf.ClientSecret)
	req.URL.RawQuery = q.Encode()

	resp, err := otelhttp.DefaultClient.Do(req)
	if err != nil {
		c.log.Printf("error calling token api at %s: %v\n", c.conf.TokenURL, err)
		return "", fmt.Errorf("token api unavailable: %v", err)
	}
	defer resp.Body.Close()

	var respBody TokenResp
	json.NewDecoder(resp.Body).Decode(&respBody)

	c.token.set(respBody.AccessToken, respBody.ExpiresIn)

	return respBody.AccessToken, nil
}
