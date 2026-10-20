package appconf_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/stretchr/testify/require"
)

func Test_Load_DotEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	const config = `DB_DSN=postgres://test:test@localhost/games
APP_HTTP_ADDRESS=localhost:8000
APP_GRPC_ADDRESS=localhost:9000
APP_DEBUG_ADDRESS=localhost:6060
APP_READTIMEOUT=5s
APP_WRITETIMEOUT=10s
JAEGER_OTLP_ENDPOINT=localhost:4318
IGDB_CLIENT_ID=test
IGDB_CLIENT_SECRET=test
IGDB_TOKEN_URL=https://example.com/token
IGDB_API_URL=https://example.com/api
IGDB_API_TIMEOUT=5s
SCHED_FETCH_IGDB_GAMES=0 0 * * *
SCHED_UPDATE_TRENDING_INDEX=0 0 * * *
SCHED_UPDATE_GAME_INFO=0 0 * * *
SCHED_PROCESS_MODERATION=0 0 * * *
REDIS_ADDR=localhost:6379
REDIS_TTL=1h
GRAYLOG_ADDR=localhost:12201
S3_REGION=test
S3_ACCESS_KEY_ID=test
S3_SECRET_ACCESS_KEY=test
S3_ENDPOINT=https://example.com
S3_BUCKET_NAME=test
S3_CDN_BASE_URL=https://cdn.example.com
S3_TIMEOUT=5s
OPENAI_API_KEY=test
OPENAI_API_URL=https://example.com
OPENAI_MODERATION_MODEL=test
OPENAI_VISION_MODEL=test
OPENAI_API_TIMEOUT=5s
LOG_LEVEL=info
AUTHAPI_ADDRESS=localhost:9001
AUTHAPI_TIMEOUT=5s
`
	for line := range strings.SplitSeq(strings.TrimSpace(config), "\n") {
		key, _, _ := strings.Cut(line, "=")
		t.Setenv(key, "")
	}
	t.Setenv("APP_HTTP_ADDRESS", "localhost:8001")
	require.NoError(t, os.WriteFile(".env", []byte(config), 0o600))

	cfg, err := appconf.Load[appconf.API]()

	require.NoError(t, err)
	require.Equal(t, "postgres://test:test@localhost/games", cfg.DB.DSN)
	require.Equal(t, "localhost:8001", cfg.Web.HTTPAddress)
	require.Equal(t, 5*time.Second, cfg.Web.ReadTimeout)
}

func Test_Load_EnvOnly(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("DB_DSN", "postgres://env:env@localhost/games")

	cfg, err := appconf.Load[appconf.Manage]()

	require.NoError(t, err)
	require.Equal(t, "postgres://env:env@localhost/games", cfg.DB.DSN)
}

func Test_Load_ReportsAllMissingKeys(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("DB_DSN", "")
	t.Setenv("APP_DEBUG_ADDRESS", "")

	_, err := appconf.Load[struct {
		DB    appconf.DB    `mapstructure:",squash"`
		Debug appconf.Debug `mapstructure:",squash"`
	}]()

	require.EqualError(t, err, "DB_DSN is required\nAPP_DEBUG_ADDRESS is required")
}

func Test_Validate_SectionWithoutValidate(t *testing.T) {
	err := appconf.Validate(struct {
		DB  appconf.DB
		Bad struct{ X string }
	}{})

	require.EqualError(t, err, "config section Bad has no Validate method")
}

func Test_Validate_AllConfigsHaveValidatableSections(t *testing.T) {
	for _, cfg := range []any{appconf.API{}, appconf.Worker{}, appconf.Manage{}} {
		err := appconf.Validate(cfg)

		require.Error(t, err)
		require.NotContains(t, err.Error(), "has no Validate method")
	}
}
