package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paranoiachains/gophkeeper-cli/internal/models"
	"github.com/paranoiachains/gophkeeper-cli/internal/server/auth/config"
)

type AuthService struct {
	db     UserStorage
	device DeviceCodeStorage
	cfg    config.AuthConfig
}

type UserStorage interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	NewUser(ctx context.Context, login, password string) (*models.User, error)
}

type DeviceCodeStorage interface {
	SaveDeviceCode(ctx context.Context, deviceCode, userCode, login string, expiresIn int64) error
	ActivateDevice(ctx context.Context, userCode string) error
	GetDeviceLoginIfActivated(ctx context.Context, deviceCode string) (string, bool, error)
}

// GetUser retrieves User from db
func (a *AuthService) GetUser(ctx context.Context, login string) (*models.User, error) {
	user, err := a.db.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// DeviceAuthorize returns device code for further token polling, user code to provide to user, expiration date in seconds
func (a *AuthService) DeviceAuthorize(ctx context.Context, login, password string) (string, string, int64, error) {
	user, err := a.db.NewUser(ctx, login, password)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create user: %w", err)
	}

	deviceCode := generateRandomCode()
	userCode := generateRandomCode()

	err = a.device.SaveDeviceCode(ctx, deviceCode, userCode, user.Login, a.cfg.DeviceCodeTTL)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to save device code: %w", err)
	}

	return deviceCode, userCode, a.cfg.DeviceCodeTTL, nil
}

// ActivateDevice finds related device code and marks it as activated
func (a *AuthService) ActivateDevice(ctx context.Context, userCode string) error {
	err := a.device.ActivateDevice(ctx, userCode)
	if err != nil {
		return fmt.Errorf("failed to activate device: %w", err)
	}
	return nil
}

// PollToken polls for token. It returns token when device_code is marked as activated.
func (a *AuthService) PollToken(ctx context.Context, deviceCode string) (string, error) {
	login, activated, err := a.device.GetDeviceLoginIfActivated(ctx, deviceCode)
	if err != nil {
		return "", fmt.Errorf("failed to get device status: %w", err)
	}
	if !activated {
		return "", fmt.Errorf("device not activated yet")
	}

	token := generateJWT(login)
	return token, nil
}

func generateRandomCode() string {
	return uuid.New().String()
}

func generateJWT(login string) string {
	return fmt.Sprintf("token-for-%s", login)
}
