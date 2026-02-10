package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
	"github.com/gophkeeper/gophkeeper/internal/domain/repository/mocks"
	"github.com/gophkeeper/gophkeeper/internal/usecase/auth"
	"go.uber.org/mock/gomock"
)

func TestRefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshToken, err := crypto.GenerateRefreshToken("user-1")
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}

	userRepo := mocks.NewMockUserRepository(ctrl)
	// Репозиторий не используется в RefreshToken

	uc := auth.NewAuthUseCase(userRepo)
	out, err := uc.RefreshToken(context.Background(), auth.RefreshTokenInput{
		RefreshToken: refreshToken,
	})

	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Error("tokens should be set")
	}
	if out.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}
}

func TestRefreshToken_RefreshTokenRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := auth.NewAuthUseCase(userRepo)

	_, err := uc.RefreshToken(context.Background(), auth.RefreshTokenInput{
		RefreshToken: "",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, auth.ErrRefreshTokenRequired) {
		t.Errorf("err = %v, want ErrRefreshTokenRequired", err)
	}
}

func TestRefreshToken_InvalidRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := auth.NewAuthUseCase(userRepo)

	_, err := uc.RefreshToken(context.Background(), auth.RefreshTokenInput{
		RefreshToken: "invalid.jwt.token",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, auth.ErrInvalidRefreshToken) {
		t.Errorf("err = %v, want ErrInvalidRefreshToken", err)
	}
}
