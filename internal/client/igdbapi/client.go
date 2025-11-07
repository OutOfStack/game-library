package igdbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	defaultTimeout = 10 * time.Second
)

const (
	gamesEndpoint     = "games"
	companiesEndpoint = "companies"

	maxLimit = 500
)

var tracer = otel.Tracer("igdbapi")

// Client represents dependencies for igdb client
type Client struct {
	log        *zap.Logger
	conf       appconf.IGDB
	token      *tokenInfo
	httpClient *http.Client
}

// New constructs Client instance
func New(log *zap.Logger, conf appconf.IGDB) (*Client, error) {
	client := &http.Client{
		Transport: observability.NewMonitoredTransport(otelhttp.NewTransport(http.DefaultTransport), "igdb"),
		Timeout:   defaultTimeout,
	}

	return &Client{
		log:        log,
		token:      &tokenInfo{},
		conf:       conf,
		httpClient: client,
	}, nil
}

// GetTopRatedGames returns top-rated games
func (c *Client) GetTopRatedGames(ctx context.Context, platformsIDs []int64, releasedBefore time.Time, minRatingsCount, minRating, limit int64) ([]TopRatedGames, error) {
	ctx, span := tracer.Start(ctx, "getTopRatedGames")
	defer span.End()

	if limit > maxLimit {
		limit = maxLimit
	}

	platformsStr := make([]string, 0, len(platformsIDs))
	for _, p := range platformsIDs {
		platformsStr = append(platformsStr, strconv.Itoa(int(p)))
	}
	platforms := strings.Join(platformsStr, ",")

	reqURL, err := url.JoinPath(c.conf.APIURL, gamesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("join url path: %v", err)
	}

	query := fmt.Sprintf(
		`fields id, cover.url, first_release_date, genres.name, name, platforms, total_rating, total_rating_count,
		slug, summary, screenshots.url, websites.type, websites.url,
		involved_companies.company.name, involved_companies.developer, involved_companies.publisher;
		sort first_release_date desc;
		where total_rating != null & total_rating_count > %d & total_rating > %d & first_release_date < %d &
		version_parent = null & parent_game = null & release_dates.platform = (%s);
		limit %d;`,
		minRatingsCount, minRating, releasedBefore.Unix(), platforms, limit)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("create get top rated games request: %v", err)
	}

	err = c.setAuthHeaders(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("set auth headers: %v", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.log.Error("request igdb api", zap.String("url", reqURL), zap.Error(err))
		return nil, fmt.Errorf("igdb api unavailable: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, rErr := io.ReadAll(resp.Body)
		if rErr != nil {
			return nil, fmt.Errorf("read response body: %v", rErr)
		}
		return nil, fmt.Errorf("%s", body)
	}

	var respBody []TopRatedGames
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, fmt.Errorf("decode response body: %v", err)
	}

	return respBody, nil
}

// GetGameInfoForUpdate returns game info for update
func (c *Client) GetGameInfoForUpdate(ctx context.Context, igdbID int64) (GameInfoForUpdate, error) {
	ctx, span := tracer.Start(ctx, "getGameInfoForUpdate")
	defer span.End()

	reqURL, err := url.JoinPath(c.conf.APIURL, gamesEndpoint)
	if err != nil {
		return GameInfoForUpdate{}, fmt.Errorf("join url path: %v", err)
	}

	query := fmt.Sprintf(
		`fields id, name, platforms, total_rating, total_rating_count, websites.type, websites.url;
		where id = %d;`,
		igdbID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(query))
	if err != nil {
		return GameInfoForUpdate{}, fmt.Errorf("create get game info for update request: %v", err)
	}

	err = c.setAuthHeaders(ctx, req)
	if err != nil {
		return GameInfoForUpdate{}, fmt.Errorf("set auth headers: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("request igdb api", zap.String("url", reqURL), zap.Error(err))
		return GameInfoForUpdate{}, fmt.Errorf("igdb api unavailable: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, rErr := io.ReadAll(resp.Body)
		if rErr != nil {
			return GameInfoForUpdate{}, fmt.Errorf("read response body: %v", rErr)
		}
		return GameInfoForUpdate{}, fmt.Errorf("%s", body)
	}

	var respBody []GameInfoForUpdate
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return GameInfoForUpdate{}, fmt.Errorf("decode response body: %v", err)
	}

	if len(respBody) == 0 {
		return GameInfoForUpdate{}, fmt.Errorf("game with igdb id %d not found in igdb api", igdbID)
	}

	return respBody[0], nil
}

// CompanyExists checks if a company with the given name exists in IGDB (case-insensitive)
func (c *Client) CompanyExists(ctx context.Context, companyName string) (bool, error) {
	ctx, span := tracer.Start(ctx, "companyExists")
	defer span.End()

	reqURL, err := url.JoinPath(c.conf.APIURL, companiesEndpoint)
	if err != nil {
		return false, fmt.Errorf("join url path: %v", err)
	}

	query := fmt.Sprintf(
		`fields id, name;
		where name ~ "%s";`,
		strings.ReplaceAll(companyName, `"`, `\"`))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(query))
	if err != nil {
		return false, fmt.Errorf("create company exists request: %v", err)
	}

	err = c.setAuthHeaders(ctx, req)
	if err != nil {
		return false, fmt.Errorf("set auth headers: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("request igdb api", zap.String("url", reqURL), zap.Error(err))
		return false, fmt.Errorf("igdb api unavailable: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, rErr := io.ReadAll(resp.Body)
		if rErr != nil {
			return false, fmt.Errorf("read response body: %v", rErr)
		}
		return false, fmt.Errorf("igdb api error: %s", body)
	}

	var respBody []CompanyInfo
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return false, fmt.Errorf("decode response body: %v", err)
	}

	return len(respBody) > 0, nil
}

// GetImageByURL downloads image by url and image type and returns data as io.ReadSeeker and file name
func (c *Client) GetImageByURL(ctx context.Context, imageURL, imageType string) (GetImageResp, error) {
	ctx, span := tracer.Start(ctx, "downloadImage")
	defer span.End()

	imageURL = getImageURL(imageURL, imageType)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return GetImageResp{}, fmt.Errorf("creating get image by url request: %v", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return GetImageResp{}, fmt.Errorf("get image by url: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.Error(err))
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetImageResp{}, fmt.Errorf("read response body: %v", err)
	}

	reader := bytes.NewReader(data)
	fileName := path.Base(request.URL.Path)
	contentType := resp.Header.Get("Content-Type")

	return GetImageResp{
		Body:        reader,
		ContentType: contentType,
		FileName:    fileName,
	}, nil
}

// returns updated image url for provided image type
func getImageURL(igdbImageURL string, imageType string) string {
	if imageType != "" {
		igdbImageURL = strings.Replace(igdbImageURL, ImageTypeThumbAlias, imageType, 1)
	}
	u, err := url.Parse(igdbImageURL)
	if err != nil {
		return ""
	}
	u.Scheme = "https"
	return u.String()
}

func (c *Client) setAuthHeaders(ctx context.Context, req *http.Request) error {
	token, err := c.accessToken(ctx)
	if err != nil {
		return fmt.Errorf("getting igdb access token: %v", err)
	}
	req.Header.Set("Client-ID", c.conf.ClientID)
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

// accessToken returns access token
func (c *Client) accessToken(ctx context.Context) (string, error) {
	ctx, span := tracer.Start(ctx, "token")
	defer span.End()

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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("calling token api", zap.String("url", c.conf.TokenURL), zap.Error(err))
		return "", fmt.Errorf("token api unavailable: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.log.Error("failed to close response body", zap.Error(err))
		}
	}()

	var respBody TokenResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("decode igdb token response: %v", err)
	}

	c.token.set(respBody.AccessToken, respBody.ExpiresIn)

	return respBody.AccessToken, nil
}
