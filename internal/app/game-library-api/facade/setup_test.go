package facade_test

import (
	"context"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/facade"
	facademock "github.com/OutOfStack/game-library/internal/app/game-library-api/facade/mocks"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	cachemock "github.com/OutOfStack/game-library/internal/pkg/cache/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TestSuite struct {
	suite.Suite
	ctx             context.Context
	ctrl            *gomock.Controller
	log             *zap.Logger
	storageMock     *facademock.MockStorage
	cacheStore      *cache.RedisStore
	redisClientMock *cachemock.MockRedisClient
	s3ClientMock    *facademock.MockS3Client
	provider        *facade.Provider
}

func (s *TestSuite) SetupTest() {
	s.ctx = s.T().Context()
	s.ctrl = gomock.NewController(s.T())
	s.storageMock = facademock.NewMockStorage(s.ctrl)
	s.log = zap.NewNop()
	s.redisClientMock = cachemock.NewMockRedisClient(s.ctrl)
	s.cacheStore = cache.NewRedisStore(s.redisClientMock, s.log)
	s.s3ClientMock = facademock.NewMockS3Client(s.ctrl)
	s.provider = facade.NewProvider(s.log, s.storageMock, s.cacheStore, s.s3ClientMock)
}

func (s *TestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
