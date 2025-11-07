package infoapi_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/api/grpc/infoapi"
	mock "github.com/OutOfStack/game-library/internal/api/grpc/infoapi/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	log            *zap.Logger
	gameFacadeMock *mock.MockGameFacade
	service        *infoapi.InfoService
}

func (s *TestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.log = zap.NewNop()
	s.gameFacadeMock = mock.NewMockGameFacade(s.ctrl)
	s.service = infoapi.NewInfoService(s.log, s.gameFacadeMock)
}

func (s *TestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
