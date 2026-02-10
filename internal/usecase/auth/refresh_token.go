package auth

import (
	"context"
	"errors"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
)

var (
	ErrRefreshTokenRequired = errors.New("refresh token is required")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

// RefreshTokenInput входные данные для обновления токена
type RefreshTokenInput struct {
	RefreshToken string
}

// RefreshTokenOutput результат обновления токена
type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// RefreshToken обновляет access токен по refresh токену
func (uc *AuthUseCase) RefreshToken(ctx context.Context, in RefreshTokenInput) (*RefreshTokenOutput, error) {
	if in.RefreshToken == "" {
		return nil, ErrRefreshTokenRequired
	}

	claims, err := crypto.ValidateToken(in.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, err := crypto.GenerateAccessToken(claims.UserID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := crypto.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(crypto.AccessTokenExpiry.Seconds()),
	}, nil
}
