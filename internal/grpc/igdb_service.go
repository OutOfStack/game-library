package grpc

import (
	"context"
	"strings"

	pb "github.com/OutOfStack/game-library/api/proto/igdb"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = otel.Tracer("grpc-igdb-service")

// IGDBAPIClient defines methods for interacting with IGDB API
type IGDBAPIClient interface {
	CompanyExists(ctx context.Context, companyName string) (bool, error)
}

// IGDBService implements the gRPC IGDBService
type IGDBService struct {
	pb.UnimplementedIGDBServiceServer
	log           *zap.Logger
	igdbAPIClient IGDBAPIClient
}

// NewIGDBService creates a new IGDBService instance
func NewIGDBService(log *zap.Logger, igdbAPIClient IGDBAPIClient) *IGDBService {
	return &IGDBService{
		log:           log,
		igdbAPIClient: igdbAPIClient,
	}
}

// CompanyExists checks if a company with the given name exists in IGDB (case-insensitive)
func (s *IGDBService) CompanyExists(ctx context.Context, req *pb.CompanyExistsRequest) (*pb.CompanyExistsResponse, error) {
	ctx, span := tracer.Start(ctx, "CompanyExists")
	defer span.End()

	// empty or whitespace-only company names return false
	if strings.TrimSpace(req.CompanyName) == "" {
		return &pb.CompanyExistsResponse{Exists: false}, nil
	}

	exists, err := s.igdbAPIClient.CompanyExists(ctx, req.CompanyName)
	if err != nil {
		// log detailed error but return generic message to client
		s.log.Error("failed to check company existence",
			zap.String("company_name", req.CompanyName),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to check company existence")
	}

	return &pb.CompanyExistsResponse{Exists: exists}, nil
}
