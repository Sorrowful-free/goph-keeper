package repository

import (
	"context"

	domainrepo "github.com/gophkeeper/gophkeeper/internal/domain/repository"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/storage"
)

// userRepo реализует domain/repository.UserRepository
type userRepo struct {
	storage *storage.Storage
}

// NewUserRepository создаёт репозиторий пользователей
func NewUserRepository(storage *storage.Storage) domainrepo.UserRepository {
	return &userRepo{storage: storage}
}

// Create создаёт пользователя
func (r *userRepo) Create(ctx context.Context, login, passwordHash string) (*models.User, error) {
	return r.storage.CreateUser(login, passwordHash)
}

// GetByLogin возвращает пользователя по логину
func (r *userRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	return r.storage.GetUserByLogin(login)
}

// GetByID возвращает пользователя по ID
func (r *userRepo) GetByID(ctx context.Context, userID string) (*models.User, error) {
	return r.storage.GetUserByID(userID)
}
