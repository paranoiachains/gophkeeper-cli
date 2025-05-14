package grpc

import (
	"context"

	"github.com/paranoiachains/gophkeeper-cli/gen/pb/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type usersServerAPI struct {
	users.UnimplementedUsersServer
	users Users
}

type Users interface {
	Login(ctx context.Context, login string, password string) (string, error)
	RegisterUser(ctx context.Context, login string, password string) (string, error)
}

func Register(gRPCServer *grpc.Server, u Users) {
	users.RegisterUsersServer(gRPCServer, &usersServerAPI{users: u})
}

func (u *usersServerAPI) Login(ctx context.Context, in *users.LoginRequest) (*users.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := u.users.Login(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &users.LoginResponse{Token: token}, nil
}

func (u *usersServerAPI) RegisterUser(ctx context.Context, in *users.RegisterRequest) (*users.RegisterResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := u.users.RegisterUser(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &users.RegisterResponse{Token: token}, nil
}
