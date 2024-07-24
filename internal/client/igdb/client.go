package igdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

const (
	gamesEndpoint = "games"

	maxLimit = 500
)

// Client represents dependencies for igdb client
type Client struct {
	log   *zap.Logger
	conf  appconf.IGDB
	token *tokenInfo
}

// New constructs Client instance
func New(log *zap.Logger, conf appconf.IGDB) (*Client, error) {
	return &Client{
		log:   log,
		token: &tokenInfo{},
		conf:  conf,
	}, nil
}

// GetTopRatedGames returns top-rated games
func (c *Client) GetTopRatedGames(ctx context.Context, platformsIDs []int64, releasedBefore time.Time, minRatingsCount, minRating, limit int64) ([]TopRatedGamesResp, error) {
	if limit > maxLimit {
		limit = maxLimit
	}

	platformsStr := make([]string, 0, len(platformsIDs))
	for _, p := range platformsIDs {
		platformsStr = append(platformsStr, strconv.Itoa(int(p)))
	}
	platforms := strings.Join(platformsStr, ",")

	reqURL, _ := url.JoinPath(c.conf.APIURL, gamesEndpoint)
	data := fmt.Sprintf(
		`fields id, cover.url, first_release_date, genres.name, name, platforms, total_rating, total_rating_count,
		slug, summary, screenshots.url, websites.category, websites.url,
		involved_companies.company.name, involved_companies.developer, involved_companies.publisher;
		sort first_release_date desc;
		where total_rating != null & total_rating_count > %d & total_rating > %d & first_release_date < %d &
		version_parent = null & parent_game = null & release_dates.platform = (%s);
		limit %d;`,
		minRatingsCount, minRating, releasedBefore.Unix(), platforms, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, fmt.Errorf("create get top rated games request: %v", err)
	}

	err = c.setAuthHeaders(ctx, &req.Header)
	if err != nil {
		return nil, fmt.Errorf("set auth headers: %v", err)
	}

	resp, err := otelhttp.DefaultClient.Do(req)
	if err != nil {
		c.log.Error("request igdb api", zap.String("url", reqURL), zap.Error(err))
		return nil, fmt.Errorf("igdb api unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, rErr := io.ReadAll(resp.Body)
		if rErr != nil {
			return nil, fmt.Errorf("read response body: %v", rErr)
		}
		return nil, fmt.Errorf("%s", body)
	}

	var respBody []TopRatedGamesResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, fmt.Errorf("decode response body: %v", err)
	}

	return respBody, nil
}

// GetImageURL returns fixed image url
func GetImageURL(igdbImageURL string, imageType string) string {
	if imageType != "" {
		igdbImageURL = strings.Replace(igdbImageURL, ImageThumbAlias, imageType, 1)
	}
	u, err := url.Parse(igdbImageURL)
	if err != nil {
		return ""
	}
	u.Scheme = "https"
	return u.String()
}

func (c *Client) setAuthHeaders(ctx context.Context, header *http.Header) error {
	token, err := c.accessToken(ctx)
	if err != nil {
		return fmt.Errorf("getting igdb access token: %v", err)
	}
	header.Set("Client-ID", c.conf.ClientID)
	header.Set("Authorization", "Bearer "+token)
	return nil
}

// accessToken returns access token
func (c *Client) accessToken(ctx context.Context) (string, error) {
	token := c.token.get()
	if token != "" {
		return token, nil
	}

	reqURL := c.conf.TokenURL
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
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
		c.log.Error("calling token api", zap.String("url", c.conf.TokenURL), zap.Error(err))
		return "", fmt.Errorf("token api unavailable: %v", err)
	}
	defer resp.Body.Close()

	var respBody TokenResp
	json.NewDecoder(resp.Body).Decode(&respBody)

	c.token.set(respBody.AccessToken, respBody.ExpiresIn)

	return respBody.AccessToken, nil
}
