package auth

import (
	"context"
	"errors"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
)

var (
	ErrLoginPasswordRequired = errors.New("login and password are required")
	ErrUserAlreadyExists     = errors.New("user already exists")
)

// RegisterUserInput входные данные для регистрации
type RegisterUserInput struct {
	Login    string
	Password string
}

// RegisterUserOutput результат регистрации
type RegisterUserOutput struct {
	UserID string
}

// RegisterUser регистрирует нового пользователя
func (uc *AuthUseCase) RegisterUser(ctx context.Context, in RegisterUserInput) (*RegisterUserOutput, error) {
	if in.Login == "" || in.Password == "" {
		return nil, ErrLoginPasswordRequired
	}

	existing, err := uc.userRepo.GetByLogin(ctx, in.Login)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := crypto.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.Create(ctx, in.Login, passwordHash)
	if err != nil {
		return nil, err
	}

	return &RegisterUserOutput{UserID: user.ID}, nil
}
