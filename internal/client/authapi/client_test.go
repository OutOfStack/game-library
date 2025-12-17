package authapi_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/client/authapi"
	authapipb "github.com/OutOfStack/game-library/pkg/proto/authapi/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type authapiServer struct {
	authapipb.UnimplementedAuthApiServiceServer
	valid bool
}

func (s *authapiServer) VerifyToken(_ context.Context, _ *authapipb.VerifyTokenRequest) (*authapipb.VerifyTokenResponse, error) {
	resp := &authapipb.VerifyTokenResponse{}
	resp.SetValid(s.valid)
	return resp, nil
}

func TestNewClientRequiresAddress(t *testing.T) {
	_, err := authapi.NewClient(authapi.Config{})
	require.Error(t, err)
}

func TestVerifyToken(t *testing.T) {
	addr, cleanup := startServer(t, true)
	defer cleanup()

	client, err := authapi.NewClient(authapi.Config{
		Address: addr,
		Timeout: time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, client.Close())
	})

	valid, err := client.VerifyToken(context.Background(), "jwt-token")
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyTokenRequiresToken(t *testing.T) {
	addr, cleanup := startServer(t, true)
	defer cleanup()

	client, err := authapi.NewClient(authapi.Config{
		Address: addr,
		Timeout: time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, client.Close())
	})

	_, err = client.VerifyToken(context.Background(), "")
	require.Error(t, err)
}

func TestCtxWithTimeout(t *testing.T) {
	t.Run("sets timeout when context has no deadline", func(t *testing.T) {
		ctx := context.Background()
		timeout := 100 * time.Millisecond

		newCtx, cancel := authapi.CtxWithTimeout(ctx, timeout)
		defer cancel()

		deadline, ok := newCtx.Deadline()
		require.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(timeout), deadline, 50*time.Millisecond)
	})

	t.Run("respects existing deadline", func(t *testing.T) {
		existingDeadline := time.Now().Add(200 * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), existingDeadline)
		defer cancel()

		newCtx, newCancel := authapi.CtxWithTimeout(ctx, 100*time.Millisecond)
		defer newCancel()

		deadline, ok := newCtx.Deadline()
		require.True(t, ok)
		assert.Equal(t, existingDeadline, deadline)
	})

	t.Run("returns working cancel function when deadline exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		newCtx, newCancel := authapi.CtxWithTimeout(ctx, 100*time.Millisecond)
		require.NotNil(t, newCancel)

		newCancel()

		select {
		case <-newCtx.Done():
		case <-time.After(100 * time.Millisecond):
			t.Fatal("context should be cancelled")
		}
	})

	t.Run("handles zero timeout", func(t *testing.T) {
		ctx := context.Background()

		newCtx, cancel := authapi.CtxWithTimeout(ctx, 0)
		defer cancel()

		_, ok := newCtx.Deadline()
		assert.False(t, ok)
	})

	t.Run("handles negative timeout", func(t *testing.T) {
		ctx := context.Background()

		newCtx, cancel := authapi.CtxWithTimeout(ctx, -1*time.Second)
		defer cancel()

		_, ok := newCtx.Deadline()
		assert.False(t, ok)
	})

	t.Run("cancel function can be called multiple times", func(_ *testing.T) {
		ctx := context.Background()

		_, cancel := authapi.CtxWithTimeout(ctx, 100*time.Millisecond)

		cancel()
		cancel()
	})
}

func startServer(t *testing.T, valid bool) (string, func()) {
	t.Helper()

	lc := &net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer()
	authapipb.RegisterAuthApiServiceServer(srv, &authapiServer{valid: valid})

	go func() {
		_ = srv.Serve(lis)
	}()

	return lis.Addr().String(), func() {
		srv.GracefulStop()
		_ = lis.Close()
	}
}
