package grpc

import (
	"context"
	"errors"
	"testing"

	pb "github.com/OutOfStack/game-library/api/proto/igdb"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run go.uber.org/mock/mockgen -source=igdb_service.go -destination=mocks/igdb_service.go -package=grpc_mock

func TestIGDBService_CompanyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockClient := NewMockIGDBAPIClient(ctrl)

	service := NewIGDBService(logger, mockClient)

	t.Run("company_exists", func(t *testing.T) {
		ctx := t.Context()
		companyName := td.String()

		mockClient.EXPECT().
			CompanyExists(gomock.Any(), companyName).
			Return(true, nil)

		req := &pb.CompanyExistsRequest{
			CompanyName: companyName,
		}

		resp, err := service.CompanyExists(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.True(t, resp.Exists)
	})

	t.Run("company_does_not_exist", func(t *testing.T) {
		ctx := t.Context()
		companyName := td.String()

		mockClient.EXPECT().
			CompanyExists(gomock.Any(), companyName).
			Return(false, nil)

		req := &pb.CompanyExistsRequest{
			CompanyName: companyName,
		}

		resp, err := service.CompanyExists(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.False(t, resp.Exists)
	})

	t.Run("empty_company_name_returns_false", func(t *testing.T) {
		ctx := t.Context()

		req := &pb.CompanyExistsRequest{
			CompanyName: "",
		}

		resp, err := service.CompanyExists(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.False(t, resp.Exists)
	})

	t.Run("internal_error", func(t *testing.T) {
		ctx := t.Context()
		companyName := td.String()

		mockClient.EXPECT().
			CompanyExists(gomock.Any(), companyName).
			Return(false, errors.New("api error"))

		req := &pb.CompanyExistsRequest{
			CompanyName: companyName,
		}

		resp, err := service.CompanyExists(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Internal, st.Code())
	})
}
