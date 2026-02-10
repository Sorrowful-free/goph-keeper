package auth

import (
	"github.com/gophkeeper/gophkeeper/internal/domain/repository"
)

// AuthUseCase объединяет сценарии аутентификации
type AuthUseCase struct {
	userRepo repository.UserRepository
}

// NewAuthUseCase создаёт use case аутентификации
func NewAuthUseCase(userRepo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}
