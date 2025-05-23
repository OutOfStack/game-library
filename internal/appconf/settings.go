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
	Auth      Auth      `mapstructure:",squash"`
	Zipkin    Zipkin    `mapstructure:",squash"`
	IGDB      IGDB      `mapstructure:",squash"`
	Scheduler Scheduler `mapstructure:",squash"`
	Redis     Redis     `mapstructure:",squash"`
	Graylog   Graylog   `mapstructure:",squash"`
	S3        S3        `mapstructure:",squash"`
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
	FetchIGDBGames string `mapstructure:"SCHED_FETCH_IGDB_GAMES"`
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

// Log represents settings for logging
type Log struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

// GetDB returns DB settings
func (cfg *Cfg) GetDB() DB {
	if cfg == nil {
		return DB{}
	}
	return cfg.DB
}

// GetWeb returns Web settings
func (cfg *Cfg) GetWeb() Web {
	if cfg == nil {
		return Web{}
	}
	return cfg.Web
}

// GetZipkin returns Zipkin settings
func (cfg *Cfg) GetZipkin() Zipkin {
	if cfg == nil {
		return Zipkin{}
	}
	return cfg.Zipkin
}

// GetAuth returns Auth settings
func (cfg *Cfg) GetAuth() Auth {
	if cfg == nil {
		return Auth{}
	}
	return cfg.Auth
}

// GetIGDB returns IGDB settings
func (cfg *Cfg) GetIGDB() IGDB {
	if cfg == nil {
		return IGDB{}
	}
	return cfg.IGDB
}

// GetScheduler returns Scheduler settings
func (cfg *Cfg) GetScheduler() Scheduler {
	if cfg == nil {
		return Scheduler{}
	}
	return cfg.Scheduler
}

// GetRedis returns Redis settings
func (cfg *Cfg) GetRedis() Redis {
	if cfg == nil {
		return Redis{}
	}
	return cfg.Redis
}

// GetS3 returns S3 settings
func (cfg *Cfg) GetS3() S3 {
	if cfg == nil {
		return S3{}
	}
	return cfg.S3
}

// GetGraylog returns Graylog settings
func (cfg *Cfg) GetGraylog() Graylog {
	if cfg == nil {
		return Graylog{}
	}
	return cfg.Graylog
}

// GetLog returns Log settings
func (cfg *Cfg) GetLog() Log {
	if cfg == nil {
		return Log{}
	}
	return cfg.Log
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

	// log
	if cfg.Log.Level == "" {
		return errors.New("LOG_LEVEL is required")
	}

	return nil
}
