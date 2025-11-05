package appconf

import (
	"errors"
	"net/url"
	"time"
)

// ServiceName - service name
const ServiceName = "game-library-api"

// Cfg - app configuration
type Cfg struct {
	Log       Log       `mapstructure:",squash"`
	DB        DB        `mapstructure:",squash"`
	Web       Web       `mapstructure:",squash"`
	GRPC      GRPC      `mapstructure:",squash"`
	Auth      Auth      `mapstructure:",squash"`
	Zipkin    Zipkin    `mapstructure:",squash"`
	IGDB      IGDB      `mapstructure:",squash"`
	Scheduler Scheduler `mapstructure:",squash"`
	Redis     Redis     `mapstructure:",squash"`
	Graylog   Graylog   `mapstructure:",squash"`
	S3        S3        `mapstructure:",squash"`
	OpenAI    OpenAI    `mapstructure:",squash"`
}

// DB represents settings for database
type DB struct {
	DSN string `mapstructure:"DB_DSN"`
}

// Web represents settings for to web server
type Web struct {
	Address           string        `mapstructure:"APP_ADDRESS"`
	DebugAddress      string        `mapstructure:"DEBUG_ADDRESS"`
	ReadTimeout       time.Duration `mapstructure:"APP_READTIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"APP_WRITETIMEOUT"`
	ShutdownTimeout   time.Duration `mapstructure:"APP_SHUTDOWNTIMEOUT"`
	AllowedCORSOrigin string        `mapstructure:"APP_ALLOWEDCORSORIGIN"`
}

// GRPC represents settings for gRPC server
type GRPC struct {
	Address string `mapstructure:"GRPC_ADDRESS"`
}

// Zipkin represents settings for Zipkin trace storage
type Zipkin struct {
	ReporterURL string `mapstructure:"ZIPKIN_REPORTERURL"`
}

// Auth represents settings for authentication and authorization
type Auth struct {
	VerifyTokenAPIURL string `mapstructure:"AUTH_VERIFYTOKENURL"`
}

// IGDB represents settings for IGDB client
type IGDB struct {
	ClientID     string `mapstructure:"IGDB_CLIENT_ID"`
	ClientSecret string `mapstructure:"IGDB_CLIENT_SECRET"`
	TokenURL     string `mapstructure:"IGDB_TOKEN_URL"`
	APIURL       string `mapstructure:"IGDB_API_URL"`
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
	Region          string `mapstructure:"S3_REGION"`
	AccessKeyID     string `mapstructure:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `mapstructure:"S3_SECRET_ACCESS_KEY"`
	Endpoint        string `mapstructure:"S3_ENDPOINT"`
	BucketName      string `mapstructure:"S3_BUCKET_NAME"`
	CDNBaseURL      string `mapstructure:"S3_CDN_BASE_URL"`
}

// OpenAI represents settings for OpenAI client
type OpenAI struct {
	APIKey          string `mapstructure:"OPENAI_API_KEY"`
	APIURL          string `mapstructure:"OPENAI_API_URL"`
	ModerationModel string `mapstructure:"OPENAI_MODERATION_MODEL"`
	VisionModel     string `mapstructure:"OPENAI_VISION_MODEL"`
}

// Log represents settings for logging
type Log struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

// Validate validates the config
func (cfg *Cfg) Validate() error {
	if cfg == nil {
		return errors.New("cfg is nil")
	}

	// db
	if cfg.DB.DSN == "" {
		return errors.New("DB_DSN is required")
	}

	// web
	if cfg.Web.Address == "" {
		return errors.New("APP_ADDRESS is required")
	}
	if cfg.Web.DebugAddress == "" {
		return errors.New("DEBUG_ADDRESS is required")
	}
	if cfg.Web.ReadTimeout <= 0 {
		return errors.New("APP_READTIMEOUT must be greater than 0")
	}
	if cfg.Web.WriteTimeout <= 0 {
		return errors.New("APP_WRITETIMEOUT must be greater than 0")
	}
	if cfg.Web.ShutdownTimeout <= 0 {
		return errors.New("APP_SHUTDOWNTIMEOUT must be greater than 0")
	}

	// grpc
	if cfg.GRPC.Address == "" {
		return errors.New("GRPC_ADDRESS is required")
	}

	// zipkin
	if cfg.Zipkin.ReporterURL == "" {
		return errors.New("ZIPKIN_REPORTERURL is required")
	}

	// auth
	if cfg.Auth.VerifyTokenAPIURL == "" {
		return errors.New("AUTH_VERIFYTOKENURL is required")
	}

	// igdb
	if cfg.IGDB.ClientID == "" {
		return errors.New("IGDB_CLIENT_ID is required")
	}
	if cfg.IGDB.ClientSecret == "" {
		return errors.New("IGDB_CLIENT_SECRET is required")
	}
	if cfg.IGDB.TokenURL == "" {
		return errors.New("IGDB_TOKEN_URL is required")
	}
	if cfg.IGDB.APIURL == "" {
		return errors.New("IGDB_API_URL is required")
	}

	// scheduler
	if cfg.Scheduler.FetchIGDBGames == "" {
		return errors.New("SCHED_FETCH_IGDB_GAMES is required")
	}
	if cfg.Scheduler.UpdateTrendingIndex == "" {
		return errors.New("SCHED_UPDATE_TRENDING_INDEX is required")
	}
	if cfg.Scheduler.UpdateGameInfo == "" {
		return errors.New("SCHED_UPDATE_GAME_INFO is required")
	}
	if cfg.Scheduler.ProcessModeration == "" {
		return errors.New("SCHED_PROCESS_MODERATION is required")
	}

	// redis
	if cfg.Redis.Address == "" {
		return errors.New("REDIS_ADDR is required")
	}
	if cfg.Redis.TTL < 0 {
		return errors.New("REDIS_TTL must be greater or equal to 0")
	}

	// graylog
	if cfg.Graylog.Address == "" {
		return errors.New("GRAYLOG_ADDR is required")
	}

	// s3
	if cfg.S3.Region == "" {
		return errors.New("S3_REGION is required")
	}
	if cfg.S3.AccessKeyID == "" {
		return errors.New("S3_ACCESS_KEY_ID is required")
	}
	if cfg.S3.SecretAccessKey == "" {
		return errors.New("S3_SECRET_ACCESS_KEY is required")
	}
	if cfg.S3.Endpoint == "" {
		return errors.New("S3_ENDPOINT is required")
	}
	if cfg.S3.BucketName == "" {
		return errors.New("S3_BUCKET_NAME is required")
	}
	if cfg.S3.CDNBaseURL == "" {
		return errors.New("S3_CDN_BASE_URL is required")
	}
	if _, err := url.Parse(cfg.S3.CDNBaseURL); err != nil {
		return errors.New("S3_CDN_BASE_URL is invalid")
	}

	// openai
	if cfg.OpenAI.APIKey == "" {
		return errors.New("OPENAI_API_KEY is required")
	}
	if cfg.OpenAI.APIURL == "" {
		return errors.New("OPENAI_API_URL is required")
	}
	if cfg.OpenAI.ModerationModel == "" {
		return errors.New("OPENAI_MODERATION_MODEL is required")
	}
	if cfg.OpenAI.VisionModel == "" {
		return errors.New("OPENAI_VISION_MODEL is required")
	}

	// log
	if cfg.Log.Level == "" {
		return errors.New("LOG_LEVEL is required")
	}

	return nil
}
