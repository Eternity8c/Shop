package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	authproto "github.com/Eternity8c/proto-shop/auth-proto/gen/go/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api  authproto.AuthServiceClient
	conn *grpc.ClientConn
	log  *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string) (*Client, error) {
	const op = "Auth.New"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cc, err := grpc.DialContext(ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := authproto.NewAuthServiceClient(cc)

	return &Client{
		api:  client,
		conn: cc,
		log:  log,
	}, nil
}

func (c *Client) Register(ctx context.Context, email, password string, fullName string) (*authproto.RegisterResponse, error) {
	const op = "Auth.Register"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.api.Register(ctx, &authproto.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	})
	if resp == nil {
		c.log.Error("resp is nil")
	}
	if err != nil {
		c.log.Error("grpc Register failed")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) Login(ctx context.Context, email, password string) (*authproto.LoginResponse, error) {
	const op = "Auth.Login"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.api.Login(ctx, &authproto.LoginRequest{
		Email:    email,
		Password: password,
	})
	if resp == nil {
		c.log.Error("resp is nil")
	}
	if err != nil {
		c.log.Error("grpc Login failed")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.api.IsAdmin(ctx, &authproto.IsAdminRequest{
		UserId: userID,
	})
	if resp == nil {
		c.log.Error("resp is nil")
	}
	if err != nil {
		c.log.Error("grpc isAdmin failed")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetIsAdmin(), nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
