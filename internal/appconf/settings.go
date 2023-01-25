package appconf

import "time"

// ServiceName - service name
const ServiceName = "game-library-api"

// Cfg - app configuration
type Cfg struct {
	DB        DB        `mapstructure:",squash"`
	Web       Web       `mapstructure:",squash"`
	Zipkin    Zipkin    `mapstructure:",squash"`
	Auth      Auth      `mapstructure:",squash"`
	IGDB      IGDB      `mapstructure:",squash"`
	Scheduler Scheduler `mapstructure:",squash"`
}

// DB represents settings related to database
type DB struct {
	Host       string `mapstructure:"DB_HOST"`
	Name       string `mapstructure:"DB_NAME"`
	User       string `mapstructure:"DB_USER"`
	Password   string `mapstructure:"DB_PASSWORD"`
	RequireSSL bool   `mapstructure:"DB_REQUIRESSL"`
}

// Web represents settings related to web server
type Web struct {
	Address           string        `mapstructure:"APP_ADDRESS"`
	DebugAddress      string        `mapstructure:"DEBUG_ADDRESS"`
	ReadTimeout       time.Duration `mapstructure:"APP_READTIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"APP_WRITETIMEOUT"`
	ShutdownTimeout   time.Duration `mapstructure:"APP_SHUTDOWNTIMEOUT"`
	AllowedCORSOrigin string        `mapstructure:"APP_ALLOWEDCORSORIGIN"`
}

// Zipkin represents settings related to zipkin trace storage
type Zipkin struct {
	ReporterURL string `mapstructure:"ZIPKIN_REPORTERURL"`
}

// Auth represents settings related to authentication and authorization
type Auth struct {
	VerifyTokenAPIURL string `mapstructure:"AUTH_VERIFYTOKENURL"`
	SigningAlgorithm  string `mapstructure:"AUTH_SIGNINGALG"`
}

// IGDB represents settings related to IGDB integration
type IGDB struct {
	ClientID     string `mapstructure:"IGDB_CLIENT_ID"`
	ClientSecret string `mapstructure:"IGDB_CLIENT_SECRET"`
	TokenURL     string `mapstructure:"IGDB_TOKEN_URL"`
	APIURL       string `mapstructure:"IGDB_API_URL"`
}

// Scheduler represents settings related to scheduler tasks
type Scheduler struct {
	FetchIGDBGames string `mapstructure:"SCHED_FETCH_IGDB_GAMES"`
}
