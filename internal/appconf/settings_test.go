package appconf_test

import (
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/stretchr/testify/require"
)

func TestCfgValidateValidConfig(t *testing.T) {
	for _, cfg := range []any{api(nil), worker(nil), &appconf.Manage{DB: validTestAPI().DB}} {
		require.NoError(t, appconf.Validate(cfg))
	}
}

func TestCfgValidateValidRedisTTLZero(t *testing.T) {
	cfg := api(func(cfg *appconf.API) {
		cfg.Redis.TTL = 0
	})

	err := appconf.Validate(cfg)

	require.NoError(t, err)
}

func TestCfgValidateErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		cfg       any
		wantError string
	}{
		{
			name: "missing db dsn",
			cfg: api(func(cfg *appconf.API) {
				cfg.DB.DSN = ""
			}),
			wantError: "DB_DSN is required",
		},
		{
			name: "missing app http address",
			cfg: api(func(cfg *appconf.API) {
				cfg.Web.HTTPAddress = ""
			}),
			wantError: "APP_HTTP_ADDRESS is required",
		},
		{
			name: "missing app debug address",
			cfg: api(func(cfg *appconf.API) {
				cfg.Debug.Address = ""
			}),
			wantError: "APP_DEBUG_ADDRESS is required",
		},
		{
			name: "missing app grpc address",
			cfg: api(func(cfg *appconf.API) {
				cfg.Web.GRPCAddress = ""
			}),
			wantError: "APP_GRPC_ADDRESS is required",
		},
		{
			name: "invalid app read timeout",
			cfg: api(func(cfg *appconf.API) {
				cfg.Web.ReadTimeout = 0
			}),
			wantError: "APP_READTIMEOUT must be greater than 0",
		},
		{
			name: "invalid app write timeout",
			cfg: api(func(cfg *appconf.API) {
				cfg.Web.WriteTimeout = 0
			}),
			wantError: "APP_WRITETIMEOUT must be greater than 0",
		},
		{
			name: "missing worker debug address",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.Debug.Address = ""
			}),
			wantError: "APP_DEBUG_ADDRESS is required",
		},
		{
			name: "missing jaeger otlp endpoint",
			cfg: api(func(cfg *appconf.API) {
				cfg.Jaeger.OTLPEndpoint = ""
			}),
			wantError: "JAEGER_OTLP_ENDPOINT is required",
		},
		{
			name: "jaeger otlp endpoint with scheme",
			cfg: api(func(cfg *appconf.API) {
				cfg.Jaeger.OTLPEndpoint = "http://jaeger:4318"
			}),
			wantError: "JAEGER_OTLP_ENDPOINT must be host:port without scheme",
		},
		{
			name: "jaeger otlp endpoint invalid format",
			cfg: api(func(cfg *appconf.API) {
				cfg.Jaeger.OTLPEndpoint = "not a valid endpoint"
			}),
			wantError: "JAEGER_OTLP_ENDPOINT must be host:port without scheme",
		},
		{
			name: "missing igdb client id",
			cfg: api(func(cfg *appconf.API) {
				cfg.IGDB.ClientID = ""
			}),
			wantError: "IGDB_CLIENT_ID is required",
		},
		{
			name: "missing igdb client secret",
			cfg: api(func(cfg *appconf.API) {
				cfg.IGDB.ClientSecret = ""
			}),
			wantError: "IGDB_CLIENT_SECRET is required",
		},
		{
			name: "missing igdb token url",
			cfg: api(func(cfg *appconf.API) {
				cfg.IGDB.TokenURL = ""
			}),
			wantError: "IGDB_TOKEN_URL is required",
		},
		{
			name: "missing igdb api url",
			cfg: api(func(cfg *appconf.API) {
				cfg.IGDB.APIURL = ""
			}),
			wantError: "IGDB_API_URL is required",
		},
		{
			name: "invalid igdb api timeout",
			cfg: api(func(cfg *appconf.API) {
				cfg.IGDB.Timeout = 0
			}),
			wantError: "IGDB_API_TIMEOUT must be greater than 0",
		},
		{
			name: "missing sched fetch igdb games",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.Scheduler.FetchIGDBGames = ""
			}),
			wantError: "SCHED_FETCH_IGDB_GAMES is required",
		},
		{
			name: "missing sched update trending index",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.Scheduler.UpdateTrendingIndex = ""
			}),
			wantError: "SCHED_UPDATE_TRENDING_INDEX is required",
		},
		{
			name: "missing sched update game info",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.Scheduler.UpdateGameInfo = ""
			}),
			wantError: "SCHED_UPDATE_GAME_INFO is required",
		},
		{
			name: "missing sched process moderation",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.Scheduler.ProcessModeration = ""
			}),
			wantError: "SCHED_PROCESS_MODERATION is required",
		},
		{
			name: "missing redis addr",
			cfg: api(func(cfg *appconf.API) {
				cfg.Redis.Address = ""
			}),
			wantError: "REDIS_ADDR is required",
		},
		{
			name: "invalid redis ttl",
			cfg: api(func(cfg *appconf.API) {
				cfg.Redis.TTL = -time.Second
			}),
			wantError: "REDIS_TTL must be greater or equal to 0",
		},
		{
			name: "missing graylog addr",
			cfg: api(func(cfg *appconf.API) {
				cfg.Graylog.Address = ""
			}),
			wantError: "GRAYLOG_ADDR is required",
		},
		{
			name: "missing s3 region",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.Region = ""
			}),
			wantError: "S3_REGION is required",
		},
		{
			name: "missing s3 access key id",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.AccessKeyID = ""
			}),
			wantError: "S3_ACCESS_KEY_ID is required",
		},
		{
			name: "missing s3 secret access key",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.SecretAccessKey = ""
			}),
			wantError: "S3_SECRET_ACCESS_KEY is required",
		},
		{
			name: "missing s3 endpoint",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.Endpoint = ""
			}),
			wantError: "S3_ENDPOINT is required",
		},
		{
			name: "missing s3 bucket name",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.BucketName = ""
			}),
			wantError: "S3_BUCKET_NAME is required",
		},
		{
			name: "missing s3 cdn base url",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.CDNBaseURL = ""
			}),
			wantError: "S3_CDN_BASE_URL is required",
		},
		{
			name: "invalid s3 cdn base url",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.CDNBaseURL = "%"
			}),
			wantError: "S3_CDN_BASE_URL is invalid",
		},
		{
			name: "invalid s3 timeout",
			cfg: api(func(cfg *appconf.API) {
				cfg.S3.Timeout = 0
			}),
			wantError: "S3_TIMEOUT must be greater than 0",
		},
		{
			name: "missing openai api key",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.OpenAI.APIKey = ""
			}),
			wantError: "OPENAI_API_KEY is required",
		},
		{
			name: "missing openai api url",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.OpenAI.APIURL = ""
			}),
			wantError: "OPENAI_API_URL is required",
		},
		{
			name: "invalid openai api timeout",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.OpenAI.Timeout = 0
			}),
			wantError: "OPENAI_API_TIMEOUT must be greater than 0",
		},
		{
			name: "missing openai moderation model",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.OpenAI.ModerationModel = ""
			}),
			wantError: "OPENAI_MODERATION_MODEL is required",
		},
		{
			name: "missing openai vision model",
			cfg: worker(func(cfg *appconf.Worker) {
				cfg.OpenAI.VisionModel = ""
			}),
			wantError: "OPENAI_VISION_MODEL is required",
		},
		{
			name: "missing log level",
			cfg: api(func(cfg *appconf.API) {
				cfg.Log.Level = ""
			}),
			wantError: "LOG_LEVEL is required",
		},
		{
			name: "missing authapi address",
			cfg: api(func(cfg *appconf.API) {
				cfg.AuthAPI.Address = ""
			}),
			wantError: "AUTHAPI_ADDRESS is required",
		},
		{
			name: "invalid authapi timeout",
			cfg: api(func(cfg *appconf.API) {
				cfg.AuthAPI.Timeout = 0
			}),
			wantError: "AUTHAPI_TIMEOUT must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := appconf.Validate(tt.cfg)

			require.EqualError(t, err, tt.wantError)
		})
	}
}

func api(mutate func(*appconf.API)) *appconf.API {
	cfg := validTestAPI()
	if mutate != nil {
		mutate(cfg)
	}
	return cfg
}

func worker(mutate func(*appconf.Worker)) *appconf.Worker {
	a := validTestAPI()
	cfg := &appconf.Worker{
		Log:    a.Log,
		DB:     a.DB,
		Debug:  a.Debug,
		Jaeger: a.Jaeger,
		IGDB:   a.IGDB,
		Scheduler: appconf.Scheduler{
			FetchIGDBGames:      "0 5 * * *",
			UpdateTrendingIndex: "0 2 * * *",
			UpdateGameInfo:      "0 1 * * *",
			ProcessModeration:   "*/2 * * * *",
		},
		Redis:   a.Redis,
		Graylog: a.Graylog,
		S3:      a.S3,
		OpenAI: appconf.OpenAI{
			APIKey:          "key",
			APIURL:          "https://api.openai.com/v1",
			ModerationModel: "omni-moderation-latest",
			VisionModel:     "gpt-6-luna",
			Timeout:         30 * time.Second,
		},
	}
	if mutate != nil {
		mutate(cfg)
	}
	return cfg
}

func validTestAPI() *appconf.API {
	return &appconf.API{
		Log: appconf.Log{
			Level: "INFO",
		},
		DB: appconf.DB{
			DSN: "postgres://user:password@localhost:5432/games?sslmode=disable",
		},
		Web: appconf.Web{
			HTTPAddress:       "0.0.0.0:8000",
			GRPCAddress:       "0.0.0.0:9000",
			ReadTimeout:       time.Second,
			WriteTimeout:      time.Second,
			AllowedCORSOrigin: "http://localhost:3000",
		},
		Debug: appconf.Debug{
			Address: "0.0.0.0:6060",
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
		AuthAPI: appconf.AuthAPI{
			Address: "localhost:9001",
			Timeout: time.Second,
		},
	}
}
