package taskprocessor_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/taskprocessor"
	mock "github.com/OutOfStack/game-library/internal/taskprocessor/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	log            *zap.Logger
	storageMock    *mock.MockStorage
	igdbClientMock *mock.MockIGDBAPIClient
	s3ClientMock   *mock.MockS3Client
	tx             *mock.MockTx
	provider       *taskprocessor.TaskProvider
}

func (s *TestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.log = zap.NewNop()
	s.storageMock = mock.NewMockStorage(s.ctrl)
	s.igdbClientMock = mock.NewMockIGDBAPIClient(s.ctrl)
	s.s3ClientMock = mock.NewMockS3Client(s.ctrl)
	s.tx = mock.NewMockTx(s.ctrl)
	s.provider = taskprocessor.New(s.log, s.storageMock, s.igdbClientMock, s.s3ClientMock)
}

func (s *TestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
