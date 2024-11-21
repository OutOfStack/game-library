package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api"
	apimock "github.com/OutOfStack/game-library/internal/app/game-library-api/api/mocks"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	cachemock "github.com/OutOfStack/game-library/internal/pkg/cache/redis/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TestSuite struct {
	suite.Suite
	Ctrl            *gomock.Controller
	GameFacadeMock  *apimock.MockGameFacade
	Log             *zap.Logger
	CacheStore      *cache.RedisStore
	RedisClientMock *cachemock.MockRedisClient
	HTTPResponse    *httptest.ResponseRecorder
	HTTPRequest     *http.Request
	Provider        *api.Provider
}

func (s *TestSuite) SetupTest() {
	s.Ctrl = gomock.NewController(s.T())
	s.GameFacadeMock = apimock.NewMockGameFacade(s.Ctrl)
	s.Log = zap.NewNop()
	s.RedisClientMock = cachemock.NewMockRedisClient(s.Ctrl)
	s.CacheStore = cache.NewRedisStore(s.RedisClientMock, s.Log)
	s.HTTPResponse = httptest.NewRecorder()
	s.HTTPRequest, _ = http.NewRequest(http.MethodGet, "/", nil)
	s.Provider = api.NewProvider(s.Log, s.CacheStore, s.GameFacadeMock)
}

func (s *TestSuite) TearDownTest() {
	s.Ctrl.Finish()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
