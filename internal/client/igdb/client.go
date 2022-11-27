package igdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/OutOfStack/game-library/internal/appconf"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	gamesEndpoint = "games"

	maxLimit = 500
)

// Client represents dependencies for igdb client
type Client struct {
	log   *log.Logger
	conf  appconf.IGDB
	token *token
}

// New constructs IGDB instance
func New(log *log.Logger, conf appconf.IGDB) (*Client, error) {
	return &Client{
		log:   log,
		token: &token{},
		conf:  conf,
	}, nil
}

// GetTopRatedGames returns top rated games
func (c *Client) GetTopRatedGames(ctx context.Context, limit uint64) ([]TopRatedGamesResp, error) {
	if limit > maxLimit {
		limit = maxLimit
	}

	reqURL, _ := url.JoinPath(c.conf.APIURL, gamesEndpoint)
	data := fmt.Sprintf(
		`fields name, rating, rating_count;
		sort rating desc;
		where rating != null & (rating_count > 100 | aggregated_rating_count > 50) & version_parent = null & parent_game = null;
		limit %d;`,
		limit)
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, fmt.Errorf("creating get top rated games request: %v", err)
	}

	err = c.setAuthHeaders(ctx, &req.Header)
	if err != nil {
		return nil, fmt.Errorf("setting auth headers: %v", err)
	}

	resp, err := otelhttp.DefaultClient.Do(req)
	if err != nil {
		c.log.Printf("error calling igdb api at %s: %v\n", reqURL, err)
		return nil, fmt.Errorf("igdb api unavailable: %v", err)
	}
	defer resp.Body.Close()

	var respBody []TopRatedGamesResp
	json.NewDecoder(resp.Body).Decode(&respBody)

	return respBody, nil
}

func (c *Client) setAuthHeaders(ctx context.Context, header *http.Header) error {
	token, err := c.accessToken(ctx)
	if err != nil {
		return fmt.Errorf("getting igdb access token: %v", err)
	}
	header.Set("Client-ID", c.conf.ClientID)
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

// accessToken returns access token
func (c *Client) accessToken(ctx context.Context) (string, error) {
	token := c.token.get()
	if token != "" {
		return token, nil
	}

	reqURL := c.conf.TokenURL
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating igdb get token request: %v", err)
	}

	q := req.URL.Query()
	q.Add("grant_type", "client_credentials")
	q.Add("client_id", c.conf.ClientID)
	q.Add("client_secret", c.conf.ClientSecret)
	req.URL.RawQuery = q.Encode()

	resp, err := otelhttp.DefaultClient.Do(req)
	if err != nil {
		c.log.Printf("error calling token api at %s: %v\n", reqURL, err)
		return "", fmt.Errorf("token api unavailable: %v", err)
	}
	defer resp.Body.Close()

	var respBody TokenResp
	json.NewDecoder(resp.Body).Decode(&respBody)

	c.token.set(respBody.AccessToken, respBody.ExpiresIn)

	return respBody.AccessToken, nil
}
