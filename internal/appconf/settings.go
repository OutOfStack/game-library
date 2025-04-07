package appconf

import "time"

// ServiceName - service name
const ServiceName = "game-library-api"

// Cfg - app configuration
type Cfg struct {
	DB         DB         `mapstructure:",squash"`
	Web        Web        `mapstructure:",squash"`
	Zipkin     Zipkin     `mapstructure:",squash"`
	Auth       Auth       `mapstructure:",squash"`
	IGDB       IGDB       `mapstructure:",squash"`
	Scheduler  Scheduler  `mapstructure:",squash"`
	Uploadcare Uploadcare `mapstructure:",squash"`
	Redis      Redis      `mapstructure:",squash"`
	Graylog    Graylog    `mapstructure:",squash"`
	S3         S3         `mapstructure:",squash"`
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
	SigningAlgorithm  string `mapstructure:"AUTH_SIGNINGALG"`
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

// Uploadcare represents settings for Uploadcare client
type Uploadcare struct {
	PublicKey string `mapstructure:"UPLOADCARE_PUBLIC_KEY"`
	SecretKey string `mapstructure:"UPLOADCARE_SECRET_KEY"`
}

// Redis represents settings for Redis client
type Redis struct {
	Address  string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	TTL      string `mapstructure:"REDIS_TTL"`
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
