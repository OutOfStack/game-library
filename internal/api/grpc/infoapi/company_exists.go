package infoapi

import (
	"strings"

	pb "github.com/OutOfStack/game-library/pkg/proto/infoapi"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CompanyExists checks if a company with the given name exists in IGDB (case-insensitive)
func (s *InfoService) CompanyExists(ctx context.Context, req *pb.CompanyExistsRequest) (*pb.CompanyExistsResponse, error) {
	ctx, span := tracer.Start(ctx, "CompanyExists")
	defer span.End()

	// empty or whitespace-only company names return false
	if strings.TrimSpace(req.GetCompanyName()) == "" {
		return nil, status.Error(codes.InvalidArgument, "empty company name")
	}

	exists, err := s.gameFacade.CompanyExistsInIGDB(ctx, req.GetCompanyName())
	if err != nil {
		// log detailed error but return generic message to client
		s.log.Error("failed to check company existence",
			zap.String("company_name", req.GetCompanyName()),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to check company existence")
	}

	return pb.CompanyExistsResponse_builder{
		Exists: proto.Bool(exists),
	}.Build(), nil
}
