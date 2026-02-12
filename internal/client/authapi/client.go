package authapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	authapipb "github.com/OutOfStack/game-library/pkg/proto/authapi/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config - settings for authapi service
type Config struct {
	Address     string
	Timeout     time.Duration
	DialOptions []grpc.DialOption
}

// Client wraps a gRPC AuthApiService client
type Client struct {
	conn    *grpc.ClientConn
	api     authapipb.AuthApiServiceClient
	timeout time.Duration
}

// NewClient dials the authapi service and returns a ready client
func NewClient(conf Config) (*Client, error) {
	if conf.Address == "" {
		return nil, errors.New("authapi address is required")
	}

	dialOpts := append([]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, conf.DialOptions...)

	conn, err := grpc.NewClient(conf.Address, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dial authapi: %w", err)
	}

	return &Client{
		timeout: conf.Timeout,
		conn:    conn,
		api:     authapipb.NewAuthApiServiceClient(conn),
	}, nil
}

// Close closes the underlying gRPC connection
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

// VerifyToken returns result of token verification
func (c *Client) VerifyToken(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, errors.New("token is required")
	}

	ctx, cancel := CtxWithTimeout(ctx, c.timeout)
	defer cancel()

	req := &authapipb.VerifyTokenRequest{}
	req.SetToken(token)

	resp, err := c.api.VerifyToken(ctx, req)
	if err != nil {
		return false, fmt.Errorf("call verify token: %w", err)
	}

	return resp.GetValid(), nil
}

// CtxWithTimeout returns context and cancel fn with provided timeout if no deadline set in context,
// otherwise returns original context and cancel fc
func CtxWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return context.WithCancel(ctx)
	}

	if timeout <= 0 {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, timeout)
}
