package grpc

import (
	"context"

	"github.com/paranoiachains/gophkeeper-cli/gen/pb/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type keeperServerAPI struct {
	keeper.UnimplementedKeeperServer
	keeper Keeper
}

type Keeper interface {
	Login(ctx context.Context, login string, password string) (string, error)
	RegisterUser(ctx context.Context, login string, password string) (string, error)
}

func Register(gRPCServer *grpc.Server, u keeper) {
	keeper.RegisterkeeperServer(gRPCServer, &keeperServerAPI{keeper: u})
}

func (u *keeperServerAPI) Login(ctx context.Context, in *keeper.LoginRequest) (*keeper.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := u.keeper.Login(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &keeper.LoginResponse{Token: token}, nil
}

func (u *keeperServerAPI) RegisterUser(ctx context.Context, in *keeper.RegisterRequest) (*keeper.RegisterResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := u.keeper.RegisterUser(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &keeper.RegisterResponse{Token: token}, nil
}
