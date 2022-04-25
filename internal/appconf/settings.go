package appconf

import "time"

// Auth represents settings related to database
type DB struct {
	Host       string `mapstructure:"DB_HOST"`
	Name       string `mapstructure:"DB_NAME"`
	User       string `mapstructure:"DB_USER"`
	Password   string `mapstructure:"DB_PASSWORD"`
	RequireSSL bool   `mapstructure:"DB_REQUIRESSL"`
}

// Auth represents settings related to web server
type Web struct {
	Address           string        `mapstructure:"APP_ADDRESS"`
	DebugAddress      string        `mapstructure:"DEBUG_ADDRESS"`
	ReadTimeout       time.Duration `mapstructure:"APP_READTIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"APP_WRITETIMEOUT"`
	ShutdownTimeout   time.Duration `mapstructure:"APP_SHUTDOWNTIMEOUT"`
	AllowedCORSOrigin string        `mapstructure:"APP_ALLOWEDCORSORIGIN"`
}

type Zipkin struct {
	ReporterURL string `mapstructure:"ZIPKIN_REPORTERURL"`
}

// Auth represents settings related to authentication and authorization
type Auth struct {
	VerifyTokenApiUrl string `mapstructure:"AUTH_VERIFYTOKENURL"`
	SigningAlgorithm  string `mapstructure:"AUTH_SIGNINGALG"`
}
