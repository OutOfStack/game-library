package appconf_test

import (
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/stretchr/testify/require"
)

func TestCfgValidateValidConfig(t *testing.T) {
	cfg := validTestCfg()

	err := cfg.Validate()

	require.NoError(t, err)
}

func TestCfgValidateValidRedisTTLZero(t *testing.T) {
	cfg := validTestCfg()
	cfg.Redis.TTL = 0

	err := cfg.Validate()

	require.NoError(t, err)
}

func TestCfgValidateErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *appconf.Cfg
		mutate    func(*appconf.Cfg)
		wantError string
	}{
		{
			name:      "nil cfg",
			cfg:       nil,
			wantError: "cfg is nil",
		},
		{
			name: "missing db dsn",
			mutate: func(cfg *appconf.Cfg) {
				cfg.DB.DSN = ""
			},
			wantError: "DB_DSN is required",
		},
		{
			name: "missing app http address",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Web.HTTPAddress = ""
			},
			wantError: "APP_HTTP_ADDRESS is required",
		},
		{
			name: "missing app debug address",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Web.DebugAddress = ""
			},
			wantError: "APP_DEBUG_ADDRESS is required",
		},
		{
			name: "missing app grpc address",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Web.GRPCAddress = ""
			},
			wantError: "APP_GRPC_ADDRESS is required",
		},
		{
			name: "invalid app read timeout",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Web.ReadTimeout = 0
			},
			wantError: "APP_READTIMEOUT must be greater than 0",
		},
		{
			name: "invalid app write timeout",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Web.WriteTimeout = 0
			},
			wantError: "APP_WRITETIMEOUT must be greater than 0",
		},
		{
			name: "missing jaeger otlp endpoint",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Jaeger.OTLPEndpoint = ""
			},
			wantError: "JAEGER_OTLP_ENDPOINT is required",
		},
		{
			name: "jaeger otlp endpoint with scheme",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Jaeger.OTLPEndpoint = "http://jaeger:4318"
			},
			wantError: "JAEGER_OTLP_ENDPOINT must be host:port without scheme",
		},
		{
			name: "jaeger otlp endpoint invalid format",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Jaeger.OTLPEndpoint = "not a valid endpoint"
			},
			wantError: "JAEGER_OTLP_ENDPOINT must be host:port without scheme",
		},
		{
			name: "missing igdb client id",
			mutate: func(cfg *appconf.Cfg) {
				cfg.IGDB.ClientID = ""
			},
			wantError: "IGDB_CLIENT_ID is required",
		},
		{
			name: "missing igdb client secret",
			mutate: func(cfg *appconf.Cfg) {
				cfg.IGDB.ClientSecret = ""
			},
			wantError: "IGDB_CLIENT_SECRET is required",
		},
		{
			name: "missing igdb token url",
			mutate: func(cfg *appconf.Cfg) {
				cfg.IGDB.TokenURL = ""
			},
			wantError: "IGDB_TOKEN_URL is required",
		},
		{
			name: "missing igdb api url",
			mutate: func(cfg *appconf.Cfg) {
				cfg.IGDB.APIURL = ""
			},
			wantError: "IGDB_API_URL is required",
		},
		{
			name: "invalid igdb api timeout",
			mutate: func(cfg *appconf.Cfg) {
				cfg.IGDB.Timeout = 0
			},
			wantError: "IGDB_API_TIMEOUT must be greater than 0",
		},
		{
			name: "missing sched fetch igdb games",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Scheduler.FetchIGDBGames = ""
			},
			wantError: "SCHED_FETCH_IGDB_GAMES is required",
		},
		{
			name: "missing sched update trending index",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Scheduler.UpdateTrendingIndex = ""
			},
			wantError: "SCHED_UPDATE_TRENDING_INDEX is required",
		},
		{
			name: "missing sched update game info",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Scheduler.UpdateGameInfo = ""
			},
			wantError: "SCHED_UPDATE_GAME_INFO is required",
		},
		{
			name: "missing sched process moderation",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Scheduler.ProcessModeration = ""
			},
			wantError: "SCHED_PROCESS_MODERATION is required",
		},
		{
			name: "missing redis addr",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Redis.Address = ""
			},
			wantError: "REDIS_ADDR is required",
		},
		{
			name: "invalid redis ttl",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Redis.TTL = -time.Second
			},
			wantError: "REDIS_TTL must be greater or equal to 0",
		},
		{
			name: "missing graylog addr",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Graylog.Address = ""
			},
			wantError: "GRAYLOG_ADDR is required",
		},
		{
			name: "missing s3 region",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.Region = ""
			},
			wantError: "S3_REGION is required",
		},
		{
			name: "missing s3 access key id",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.AccessKeyID = ""
			},
			wantError: "S3_ACCESS_KEY_ID is required",
		},
		{
			name: "missing s3 secret access key",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.SecretAccessKey = ""
			},
			wantError: "S3_SECRET_ACCESS_KEY is required",
		},
		{
			name: "missing s3 endpoint",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.Endpoint = ""
			},
			wantError: "S3_ENDPOINT is required",
		},
		{
			name: "missing s3 bucket name",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.BucketName = ""
			},
			wantError: "S3_BUCKET_NAME is required",
		},
		{
			name: "missing s3 cdn base url",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.CDNBaseURL = ""
			},
			wantError: "S3_CDN_BASE_URL is required",
		},
		{
			name: "invalid s3 cdn base url",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.CDNBaseURL = "%"
			},
			wantError: "S3_CDN_BASE_URL is invalid",
		},
		{
			name: "invalid s3 timeout",
			mutate: func(cfg *appconf.Cfg) {
				cfg.S3.Timeout = 0
			},
			wantError: "S3_TIMEOUT must be greater than 0",
		},
		{
			name: "missing openai api key",
			mutate: func(cfg *appconf.Cfg) {
				cfg.OpenAI.APIKey = ""
			},
			wantError: "OPENAI_API_KEY is required",
		},
		{
			name: "missing openai api url",
			mutate: func(cfg *appconf.Cfg) {
				cfg.OpenAI.APIURL = ""
			},
			wantError: "OPENAI_API_URL is required",
		},
		{
			name: "invalid openai api timeout",
			mutate: func(cfg *appconf.Cfg) {
				cfg.OpenAI.Timeout = 0
			},
			wantError: "OPENAI_API_TIMEOUT must be greater than 0",
		},
		{
			name: "missing openai moderation model",
			mutate: func(cfg *appconf.Cfg) {
				cfg.OpenAI.ModerationModel = ""
			},
			wantError: "OPENAI_MODERATION_MODEL is required",
		},
		{
			name: "missing openai vision model",
			mutate: func(cfg *appconf.Cfg) {
				cfg.OpenAI.VisionModel = ""
			},
			wantError: "OPENAI_VISION_MODEL is required",
		},
		{
			name: "missing log level",
			mutate: func(cfg *appconf.Cfg) {
				cfg.Log.Level = ""
			},
			wantError: "LOG_LEVEL is required",
		},
		{
			name: "missing authapi address",
			mutate: func(cfg *appconf.Cfg) {
				cfg.AuthAPI.Address = ""
			},
			wantError: "AUTHAPI_ADDRESS is required",
		},
		{
			name: "invalid authapi timeout",
			mutate: func(cfg *appconf.Cfg) {
				cfg.AuthAPI.Timeout = 0
			},
			wantError: "AUTHAPI_TIMEOUT must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg
			if cfg == nil && tt.mutate != nil {
				cfg = validTestCfg()
			}
			if tt.mutate != nil {
				tt.mutate(cfg)
			}

			err := cfg.Validate()

			require.EqualError(t, err, tt.wantError)
		})
	}
}

func validTestCfg() *appconf.Cfg {
	return &appconf.Cfg{
		Log: appconf.Log{
			Level: "INFO",
		},
		DB: appconf.DB{
			DSN: "postgres://user:password@localhost:5432/games?sslmode=disable",
		},
		Web: appconf.Web{
			HTTPAddress:       "0.0.0.0:8000",
			GRPCAddress:       "0.0.0.0:9000",
			DebugAddress:      "0.0.0.0:6060",
			ReadTimeout:       time.Second,
			WriteTimeout:      time.Second,
			AllowedCORSOrigin: "http://localhost:3000",
		},
		Jaeger: appconf.Jaeger{
			OTLPEndpoint: "jaeger-service:4318",
		},
		IGDB: appconf.IGDB{
			ClientID:     "id",
			ClientSecret: "secret",
			TokenURL:     "https://id.twitch.tv/oauth2/token",
			APIURL:       "https://api.igdb.com/v4/",
			Timeout:      10 * time.Second,
		},
		Scheduler: appconf.Scheduler{
			FetchIGDBGames:      "0 5 * * *",
			UpdateTrendingIndex: "0 2 * * *",
			UpdateGameInfo:      "0 1 * * *",
			ProcessModeration:   "*/2 * * * *",
		},
		Redis: appconf.Redis{
			Address:  "localhost:6379",
			Password: "",
			TTL:      2 * time.Hour,
		},
		Graylog: appconf.Graylog{
			Address: "localhost:12201",
		},
		S3: appconf.S3{
			Region:          "auto",
			AccessKeyID:     "access-key",
			SecretAccessKey: "secret-key",
			Endpoint:        "https://s3.example.com",
			BucketName:      "game-library",
			CDNBaseURL:      "https://cdn.example.com",
			Timeout:         10 * time.Second,
		},
		OpenAI: appconf.OpenAI{
			APIKey:          "key",
			APIURL:          "https://api.openai.com/v1",
			ModerationModel: "omni-moderation-latest",
			VisionModel:     "gpt-5-nano",
			Timeout:         30 * time.Second,
		},
		AuthAPI: appconf.AuthAPI{
			Address: "localhost:9001",
			Timeout: time.Second,
		},
	}
}
