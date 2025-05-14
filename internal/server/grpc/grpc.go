package grpc

import (
	"context"

	"github.com/paranoiachains/gophkeeper-cli/gen/pb"
	"github.com/paranoiachains/gophkeeper-cli/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type usersServerAPI struct {
	pb.UnimplementedUsersServer
	users Users
}

type Users interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	Login(ctx context.Context, login string, password string) (string, error)
	RegisterUser(ctx context.Context, login string, password string) (string, error)
}

func Register(gRPCServer *grpc.Server, users Users) {
	pb.RegisterUsersServer(gRPCServer, &usersServerAPI{users: users})
}

func (u *usersServerAPI) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
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

	return &pb.LoginResponse{Token: token}, nil
}

func (u *usersServerAPI) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required in order to get user")
	}

	user, err := u.users.GetUser(ctx, in.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{User: &pb.User{
		Id:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}}, nil
}

func (u *usersServerAPI) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
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

	return &pb.RegisterUserResponse{Token: token}, nil
}
