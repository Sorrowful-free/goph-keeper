package auth

import (
	"context"
	"errors"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
)

// LoginUserInput входные данные для входа
type LoginUserInput struct {
	Login    string
	Password string
}

// LoginUserOutput результат входа (токены)
type LoginUserOutput struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// LoginUser выполняет вход пользователя и возвращает токены
func (uc *AuthUseCase) LoginUser(ctx context.Context, in LoginUserInput) (*LoginUserOutput, error) {
	if in.Login == "" || in.Password == "" {
		return nil, ErrLoginPasswordRequired
	}

	user, err := uc.userRepo.GetByLogin(ctx, in.Login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !crypto.CheckPassword(in.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := crypto.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := crypto.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(crypto.AccessTokenExpiry.Seconds()),
	}, nil
}
