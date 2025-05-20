package auth_test

import (
	"context"
	"testing"

	"github.com/OutOfStack/game-library/internal/auth"
	mock "github.com/OutOfStack/game-library/internal/auth/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	ctx           context.Context
	log           *zap.Logger
	authAPIClient *mock.MockAPIClient
	auth          *auth.Client
}

func (s *TestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ctx = context.Background()
	s.log = zap.NewNop()
	s.authAPIClient = mock.NewMockAPIClient(s.ctrl)
	var err error
	s.auth, err = auth.New(s.log, s.authAPIClient)
	s.NoError(err)
}

func (s *TestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
