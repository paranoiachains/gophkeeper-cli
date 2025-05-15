package auth

import (
	"context"

	"github.com/paranoiachains/gophkeeper-cli/gen/pb/auth"
	"github.com/paranoiachains/gophkeeper-cli/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Implements AuthClient interface
type Client struct {
	conn   *grpc.ClientConn
	client auth.AuthClient
}

// Returns *Client, uses provided address to connect to gRPC server
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := auth.NewAuthClient(conn)
	return &Client{conn: conn, client: client}, nil
}

func (c *Client) GetUser(ctx context.Context, login string) (*models.User, error) {
	resp, err := c.client.GetUser(ctx, &auth.GetUserRequest{Login: login})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       resp.User.Id,
		Login:    resp.User.Login,
		Password: resp.User.Password,
	}, nil
}

func (c *Client) DeviceAuthorize(ctx context.Context, login string) (string, string, int64, error) {
	resp, err := c.client.DeviceAuthorize(ctx, &auth.DeviceAuthorizeRequest{Login: login})
	if err != nil {
		return "", "", 0, err
	}

	return resp.DeviceCode, resp.UserCode, resp.ExpiresIn, nil
}

func (c *Client) PollToken(ctx context.Context, deviceCode string) (*auth.PollTokenResponse, error) {
	resp, err := c.client.PollToken(ctx, &auth.PollTokenRequest{
		DeviceCode: deviceCode,
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}
