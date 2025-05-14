package service

import (
	"context"

	"github.com/paranoiachains/gophkeeper-cli/internal/models"
)

// Implementations of Users interface
type Users struct {
	db UserStorage
}

// Users DB layer
type UserStorage interface {
	GetUser(ctx context.Context, login string) (*models.User, error)

	// Returns user ID for token building
	Login(ctx context.Context, login string, password string) (string, error)

	// Also returns user ID for token building
	RegisterUser(ctx context.Context, login string, password string) (string, error)
}
