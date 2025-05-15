package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrAuthorizationPending = errors.New("pending authorization")

// Implementations of Users interface
type Users struct {
	auth AuthClient
}

// Returns an instance of Users struct, which implements Users interface
func NewUsers(auth AuthClient) *Users {
	return &Users{auth: auth}
}

type AuthClient interface {
	DeviceAuthorize(ctx context.Context, login, password string) (deviceCode, userCode string, expiresIn int64, err error)
	ActivateDevice(ctx context.Context, userCode string) error
	PollToken(ctx context.Context, deviceCode string) (token string, err error)
}

func (u *Users) Login(ctx context.Context, login string, password string) (string, error) {
	deviceCode, userCode, expiresIn, err := u.auth.DeviceAuthorize(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("device authorization failed: %w", err)
	}

	fmt.Printf("Please enter the following code on your device to authorize: %s\n", userCode)

	timeout := time.After(time.Duration(expiresIn) * time.Second)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return "", errors.New("authorization timeout")
		case <-tick:
			token, err := u.auth.PollToken(ctx, deviceCode)
			if err == nil {
				return token, nil
			}
			// If not authorized yet, continue polling
			// If error is permanent — break
			if !errors.Is(err, ErrAuthorizationPending) {
				return "", fmt.Errorf("polling error: %w", err)
			}
		}
	}
}

func (u *Users) RegisterUser(ctx context.Context, login string, password string) (string, error) {
	return u.Login(ctx, login, password)
}
