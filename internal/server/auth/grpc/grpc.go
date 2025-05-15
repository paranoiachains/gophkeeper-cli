package grpc

import (
	"context"

	"github.com/paranoiachains/gophkeeper-cli/gen/pb/auth"
	"github.com/paranoiachains/gophkeeper-cli/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServerAPI struct {
	auth.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	DeviceAuthorize(ctx context.Context, login string) (deviceCode, userCode string, expiresIn int64, err error)
	PollToken(ctx context.Context, deviceCode string) (token string, err error)
}

func RegisterAuth(gRPCServer *grpc.Server, a Auth) {
	auth.RegisterAuthServer(gRPCServer, &authServerAPI{auth: a})
}

func (a *authServerAPI) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	user, err := a.auth.GetUser(ctx, req.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &auth.GetUserResponse{User: &auth.User{
		Id:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}}, nil
}

func (a *authServerAPI) DeviceAuthorize(ctx context.Context, req *auth.DeviceAuthorizeRequest) (*auth.DeviceAuthorizeResponse, error) {
	deviceCode, userCode, expiresIn, err := a.auth.DeviceAuthorize(ctx, req.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to authorize device")
	}

	return &auth.DeviceAuthorizeResponse{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		ExpiresIn:  expiresIn,
	}, nil
}

func (a *authServerAPI) PollToken(ctx context.Context, req *auth.PollTokenRequest) (*auth.PollTokenResponse, error) {
	token, err := a.auth.PollToken(ctx, req.DeviceCode)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to poll token")
	}
	return &auth.PollTokenResponse{
		Token: token,
	}, nil
}
