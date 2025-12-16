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

const defaultTimeout = 5 * time.Second

// Config - settings for infoapi service
type Config struct {
	Address     string
	Timeout     time.Duration
	DialOptions []grpc.DialOption
}

// Client wraps a gRPC InfoApiService client
type Client struct {
	cfg  Config
	conn *grpc.ClientConn
	api  authapipb.AuthApiServiceClient
}

// NewClient dials the infoapi service and returns a ready client
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, errors.New("infoapi address is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}

	dialOpts := append([]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, cfg.DialOptions...)

	conn, err := grpc.NewClient(cfg.Address, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dial infoapi: %w", err)
	}

	return &Client{
		cfg:  cfg,
		conn: conn,
		api:  authapipb.NewAuthApiServiceClient(conn),
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

	ctx, cancel := CtxWithTimeout(ctx, c.cfg.Timeout)
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
