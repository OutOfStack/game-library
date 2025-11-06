package infoapi

import (
	"context"

	infopb "github.com/OutOfStack/game-library/pkg/proto/infoapi"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("grpc-info-service")

// GameFacade defines methods for interacting game facade
type GameFacade interface {
	CompanyExistsInIGDB(ctx context.Context, companyName string) (bool, error)
}

// InfoService implements the gRPC InfoService
type InfoService struct {
	infopb.UnimplementedInfoApiServiceServer
	log        *zap.Logger
	gameFacade GameFacade
}

// NewInfoService creates a new InfoService instance
func NewInfoService(log *zap.Logger, gameFacade GameFacade) *InfoService {
	return &InfoService{
		log:        log,
		gameFacade: gameFacade,
	}
}
