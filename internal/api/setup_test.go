package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OutOfStack/game-library/internal/api"
	apimock "github.com/OutOfStack/game-library/internal/api/mocks"
	"github.com/OutOfStack/game-library/internal/appconf"
	mwmock "github.com/OutOfStack/game-library/internal/middleware/mocks"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	cachemock "github.com/OutOfStack/game-library/internal/pkg/cache/mocks"
	"github.com/OutOfStack/game-library/internal/web"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	gameFacadeMock  *apimock.MockGameFacade
	log             *zap.Logger
	cacheStore      *cache.RedisStore
	redisClientMock *cachemock.MockRedisClient
	authClientMock  *mwmock.MockAuthClient
	httpResponse    *httptest.ResponseRecorder
	httpRequest     *http.Request
	provider        *api.Provider
}

func (s *TestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.gameFacadeMock = apimock.NewMockGameFacade(s.ctrl)
	s.log = zap.NewNop()
	s.redisClientMock = cachemock.NewMockRedisClient(s.ctrl)
	s.cacheStore = cache.NewRedisStore(s.redisClientMock, s.log)
	s.authClientMock = mwmock.NewMockAuthClient(s.ctrl)
	s.httpResponse = httptest.NewRecorder()
	s.httpRequest, _ = http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/", nil)
	s.provider = api.NewProvider(s.log, s.cacheStore, s.gameFacadeMock, web.NewDecoder(s.log, &appconf.Cfg{}))
}

func (s *TestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
