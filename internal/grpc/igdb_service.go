package grpc

import (
	"context"
	"fmt"

	pb "github.com/OutOfStack/game-library/api/proto/igdb"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = otel.Tracer("grpc-igdb-service")

// IGDBAPIClient defines methods for interacting with IGDB API
type IGDBAPIClient interface {
	GetGameInfoForUpdate(ctx context.Context, igdbID int64) (igdbapi.GameInfoForUpdate, error)
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

// GetGameInfoForUpdate retrieves game info needed for updates from IGDB API
func (s *IGDBService) GetGameInfoForUpdate(ctx context.Context, req *pb.GetGameInfoForUpdateRequest) (*pb.GetGameInfoForUpdateResponse, error) {
	ctx, span := tracer.Start(ctx, "GetGameInfoForUpdate")
	defer span.End()

	if req.IgdbId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "igdb_id must be greater than 0")
	}

	gameInfo, err := s.igdbAPIClient.GetGameInfoForUpdate(ctx, req.IgdbId)
	if err != nil {
		s.log.Error("failed to get game info for update",
			zap.Int64("igdb_id", req.IgdbId),
			zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get game info: %v", err))
	}

	// convert websites
	websites := make([]*pb.Website, 0, len(gameInfo.Websites))
	for _, w := range gameInfo.Websites {
		websites = append(websites, &pb.Website{
			Url:  w.URL,
			Type: int32(w.Type),
		})
	}

	return &pb.GetGameInfoForUpdateResponse{
		Id:               gameInfo.ID,
		Name:             gameInfo.Name,
		TotalRating:      gameInfo.TotalRating,
		TotalRatingCount: gameInfo.TotalRatingCount,
		Platforms:        gameInfo.Platforms,
		Websites:         websites,
	}, nil
}
