package appconf

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// ServiceName - service name
const ServiceName = "game-library-api"

// API - api service configuration
type API struct {
	Log     Log     `mapstructure:",squash"`
	DB      DB      `mapstructure:",squash"`
	Web     Web     `mapstructure:",squash"`
	Debug   Debug   `mapstructure:",squash"`
	Jaeger  Jaeger  `mapstructure:",squash"`
	IGDB    IGDB    `mapstructure:",squash"`
	Redis   Redis   `mapstructure:",squash"`
	Graylog Graylog `mapstructure:",squash"`
	S3      S3      `mapstructure:",squash"`
	AuthAPI AuthAPI `mapstructure:",squash"`
}

// Worker - background worker configuration
type Worker struct {
	Log       Log       `mapstructure:",squash"`
	DB        DB        `mapstructure:",squash"`
	Debug     Debug     `mapstructure:",squash"`
	Jaeger    Jaeger    `mapstructure:",squash"`
	IGDB      IGDB      `mapstructure:",squash"`
	Scheduler Scheduler `mapstructure:",squash"`
	Redis     Redis     `mapstructure:",squash"`
	Graylog   Graylog   `mapstructure:",squash"`
	S3        S3        `mapstructure:",squash"`
	OpenAI    OpenAI    `mapstructure:",squash"`
}

// Manage - database management tool configuration
type Manage struct {
	DB DB `mapstructure:",squash"`
}

// DB represents settings for database
type DB struct {
	DSN string `mapstructure:"DB_DSN"`
}

// Web represents settings for to web server
type Web struct {
	HTTPAddress       string        `mapstructure:"APP_HTTP_ADDRESS"`
	GRPCAddress       string        `mapstructure:"APP_GRPC_ADDRESS"`
	ReadTimeout       time.Duration `mapstructure:"APP_READTIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"APP_WRITETIMEOUT"`
	AllowedCORSOrigin string        `mapstructure:"APP_ALLOWEDCORSORIGIN"`
}

// Debug represents settings for debug server (pprof, metrics, health probes)
type Debug struct {
	Address string `mapstructure:"APP_DEBUG_ADDRESS"`
}

// Jaeger represents settings for Jaeger OTLP trace export
type Jaeger struct {
	OTLPEndpoint string `mapstructure:"JAEGER_OTLP_ENDPOINT"`
}

// IGDB represents settings for IGDB client
type IGDB struct {
	ClientID     string        `mapstructure:"IGDB_CLIENT_ID"`
	ClientSecret string        `mapstructure:"IGDB_CLIENT_SECRET"`
	TokenURL     string        `mapstructure:"IGDB_TOKEN_URL"`
	APIURL       string        `mapstructure:"IGDB_API_URL"`
	Timeout      time.Duration `mapstructure:"IGDB_API_TIMEOUT"`
}

// Scheduler represents settings for task scheduler
type Scheduler struct {
	FetchIGDBGames      string `mapstructure:"SCHED_FETCH_IGDB_GAMES"`
	UpdateTrendingIndex string `mapstructure:"SCHED_UPDATE_TRENDING_INDEX"`
	UpdateGameInfo      string `mapstructure:"SCHED_UPDATE_GAME_INFO"`
	ProcessModeration   string `mapstructure:"SCHED_PROCESS_MODERATION"`
}

// Redis represents settings for Redis client
type Redis struct {
	Address  string        `mapstructure:"REDIS_ADDR"`
	Password string        `mapstructure:"REDIS_PASSWORD"`
	TTL      time.Duration `mapstructure:"REDIS_TTL"`
}

// Graylog represents settings related to Graylog integration
type Graylog struct {
	Address string `mapstructure:"GRAYLOG_ADDR"`
}

// S3 represents settings for S3 client
type S3 struct {
	Region          string        `mapstructure:"S3_REGION"`
	AccessKeyID     string        `mapstructure:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string        `mapstructure:"S3_SECRET_ACCESS_KEY"`
	Endpoint        string        `mapstructure:"S3_ENDPOINT"`
	BucketName      string        `mapstructure:"S3_BUCKET_NAME"`
	CDNBaseURL      string        `mapstructure:"S3_CDN_BASE_URL"`
	Timeout         time.Duration `mapstructure:"S3_TIMEOUT"`
}

// OpenAI represents settings for OpenAI client
type OpenAI struct {
	APIKey          string        `mapstructure:"OPENAI_API_KEY"`
	APIURL          string        `mapstructure:"OPENAI_API_URL"`
	ModerationModel string        `mapstructure:"OPENAI_MODERATION_MODEL"`
	VisionModel     string        `mapstructure:"OPENAI_VISION_MODEL"`
	Timeout         time.Duration `mapstructure:"OPENAI_API_TIMEOUT"`
}

// Log represents settings for logging
type Log struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

// AuthAPI represents settings for AuthAPI gRPC client
type AuthAPI struct {
	Address string        `mapstructure:"AUTHAPI_ADDRESS"`
	Timeout time.Duration `mapstructure:"AUTHAPI_TIMEOUT"`
}

// Validate validates db settings
func (c DB) Validate() error {
	if c.DSN == "" {
		return errors.New("DB_DSN is required")
	}
	return nil
}

// Validate validates web server settings
func (c Web) Validate() error {
	if c.HTTPAddress == "" {
		return errors.New("APP_HTTP_ADDRESS is required")
	}
	if c.GRPCAddress == "" {
		return errors.New("APP_GRPC_ADDRESS is required")
	}
	if c.ReadTimeout <= 0 {
		return errors.New("APP_READTIMEOUT must be greater than 0")
	}
	if c.WriteTimeout <= 0 {
		return errors.New("APP_WRITETIMEOUT must be greater than 0")
	}
	return nil
}

// Validate validates debug server settings
func (c Debug) Validate() error {
	if c.Address == "" {
		return errors.New("APP_DEBUG_ADDRESS is required")
	}
	return nil
}

// Validate validates Jaeger settings
func (c Jaeger) Validate() error {
	if c.OTLPEndpoint == "" {
		return errors.New("JAEGER_OTLP_ENDPOINT is required")
	}
	if strings.Contains(c.OTLPEndpoint, "://") {
		return errors.New("JAEGER_OTLP_ENDPOINT must be host:port without scheme")
	}
	u, err := url.Parse("http://" + c.OTLPEndpoint)
	if err != nil || u.Host != c.OTLPEndpoint {
		return errors.New("JAEGER_OTLP_ENDPOINT must be host:port without scheme")
	}
	return nil
}

// Validate validates IGDB settings
func (c IGDB) Validate() error {
	if c.ClientID == "" {
		return errors.New("IGDB_CLIENT_ID is required")
	}
	if c.ClientSecret == "" {
		return errors.New("IGDB_CLIENT_SECRET is required")
	}
	if c.TokenURL == "" {
		return errors.New("IGDB_TOKEN_URL is required")
	}
	if c.APIURL == "" {
		return errors.New("IGDB_API_URL is required")
	}
	if c.Timeout <= 0 {
		return errors.New("IGDB_API_TIMEOUT must be greater than 0")
	}
	return nil
}

// Validate validates scheduler settings
func (c Scheduler) Validate() error {
	if c.FetchIGDBGames == "" {
		return errors.New("SCHED_FETCH_IGDB_GAMES is required")
	}
	if c.UpdateTrendingIndex == "" {
		return errors.New("SCHED_UPDATE_TRENDING_INDEX is required")
	}
	if c.UpdateGameInfo == "" {
		return errors.New("SCHED_UPDATE_GAME_INFO is required")
	}
	if c.ProcessModeration == "" {
		return errors.New("SCHED_PROCESS_MODERATION is required")
	}
	return nil
}

// Validate validates Redis settings
func (c Redis) Validate() error {
	if c.Address == "" {
		return errors.New("REDIS_ADDR is required")
	}
	if c.TTL < 0 {
		return errors.New("REDIS_TTL must be greater or equal to 0")
	}
	return nil
}

// Validate validates Graylog settings
func (c Graylog) Validate() error {
	if c.Address == "" {
		return errors.New("GRAYLOG_ADDR is required")
	}
	return nil
}

// Validate validates S3 settings
func (c S3) Validate() error {
	if c.Region == "" {
		return errors.New("S3_REGION is required")
	}
	if c.AccessKeyID == "" {
		return errors.New("S3_ACCESS_KEY_ID is required")
	}
	if c.SecretAccessKey == "" {
		return errors.New("S3_SECRET_ACCESS_KEY is required")
	}
	if c.Endpoint == "" {
		return errors.New("S3_ENDPOINT is required")
	}
	if c.BucketName == "" {
		return errors.New("S3_BUCKET_NAME is required")
	}
	if c.CDNBaseURL == "" {
		return errors.New("S3_CDN_BASE_URL is required")
	}
	if _, err := url.Parse(c.CDNBaseURL); err != nil {
		return errors.New("S3_CDN_BASE_URL is invalid")
	}
	if c.Timeout <= 0 {
		return errors.New("S3_TIMEOUT must be greater than 0")
	}
	return nil
}

// Validate validates OpenAI settings
func (c OpenAI) Validate() error {
	if c.APIKey == "" {
		return errors.New("OPENAI_API_KEY is required")
	}
	if c.APIURL == "" {
		return errors.New("OPENAI_API_URL is required")
	}
	if c.ModerationModel == "" {
		return errors.New("OPENAI_MODERATION_MODEL is required")
	}
	if c.VisionModel == "" {
		return errors.New("OPENAI_VISION_MODEL is required")
	}
	if c.Timeout <= 0 {
		return errors.New("OPENAI_API_TIMEOUT must be greater than 0")
	}
	return nil
}

// Validate validates log settings
func (c Log) Validate() error {
	if c.Level == "" {
		return errors.New("LOG_LEVEL is required")
	}
	return nil
}

// Validate validates AuthAPI settings
func (c AuthAPI) Validate() error {
	if c.Address == "" {
		return errors.New("AUTHAPI_ADDRESS is required")
	}
	if c.Timeout <= 0 {
		return errors.New("AUTHAPI_TIMEOUT must be greater than 0")
	}
	return nil
}
