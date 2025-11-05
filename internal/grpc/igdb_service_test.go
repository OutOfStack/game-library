package grpc

import (
	"context"
	"errors"
	"testing"

	pb "github.com/OutOfStack/game-library/api/proto/igdb"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run go.uber.org/mock/mockgen -source=igdb_service.go -destination=mocks/igdb_service.go -package=grpc_mock

func TestIGDBService_GetGameInfoForUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockClient := NewMockIGDBAPIClient(ctrl)

	service := NewIGDBService(logger, mockClient)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		igdbID := td.Int64()
		if igdbID <= 0 {
			igdbID = 12345
		}

		gameInfo := igdbapi.GameInfoForUpdate{
			ID:               igdbID,
			Name:             td.String(),
			TotalRating:      td.Float64n(100),
			TotalRatingCount: td.Int32(),
			Platforms:        []int64{td.Int64(), td.Int64()},
			Websites: []igdbapi.Website{
				{URL: td.URL(), Type: int8(td.Intn(20))},
				{URL: td.URL(), Type: int8(td.Intn(20))},
			},
		}

		mockClient.EXPECT().
			GetGameInfoForUpdate(gomock.Any(), igdbID).
			Return(gameInfo, nil)

		req := &pb.GetGameInfoForUpdateRequest{
			IgdbId: igdbID,
		}

		resp, err := service.GetGameInfoForUpdate(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, gameInfo.ID, resp.Id)
		require.Equal(t, gameInfo.Name, resp.Name)
		require.Equal(t, gameInfo.TotalRating, resp.TotalRating)
		require.Equal(t, gameInfo.TotalRatingCount, resp.TotalRatingCount)
		require.Equal(t, gameInfo.Platforms, resp.Platforms)
		require.Len(t, resp.Websites, len(gameInfo.Websites))
		for i, w := range gameInfo.Websites {
			require.Equal(t, w.URL, resp.Websites[i].Url)
			require.Equal(t, int32(w.Type), resp.Websites[i].Type)
		}
	})

	t.Run("invalid_argument_zero_id", func(t *testing.T) {
		ctx := context.Background()

		req := &pb.GetGameInfoForUpdateRequest{
			IgdbId: 0,
		}

		resp, err := service.GetGameInfoForUpdate(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("invalid_argument_negative_id", func(t *testing.T) {
		ctx := context.Background()

		req := &pb.GetGameInfoForUpdateRequest{
			IgdbId: -1,
		}

		resp, err := service.GetGameInfoForUpdate(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("internal_error", func(t *testing.T) {
		ctx := context.Background()
		igdbID := int64(12345)

		mockClient.EXPECT().
			GetGameInfoForUpdate(gomock.Any(), igdbID).
			Return(igdbapi.GameInfoForUpdate{}, errors.New("api error"))

		req := &pb.GetGameInfoForUpdateRequest{
			IgdbId: igdbID,
		}

		resp, err := service.GetGameInfoForUpdate(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Internal, st.Code())
	})
}
