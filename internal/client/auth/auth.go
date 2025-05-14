package auth

import (
	"github.com/paranoiachains/gophkeeper-cli/gen/pb/auth"
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
