package repository

import (
	"context"

	"github.com/gophkeeper/gophkeeper/internal/models"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_user_repository.go -package=mocks github.com/gophkeeper/gophkeeper/internal/domain/repository UserRepository

// UserRepository определяет контракт для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, login, passwordHash string) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	GetByID(ctx context.Context, userID string) (*models.User, error)
}
