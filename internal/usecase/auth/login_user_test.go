package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
	"github.com/gophkeeper/gophkeeper/internal/domain/repository/mocks"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/usecase/auth"
	"go.uber.org/mock/gomock"
)

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, err := crypto.HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	user := &models.User{ID: "user-1", Login: "testuser", PasswordHash: hash}

	userRepo := mocks.NewMockUserRepository(ctrl)
	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	uc := auth.NewAuthUseCase(userRepo)
	out, err := uc.LoginUser(context.Background(), auth.LoginUserInput{
		Login:    "testuser",
		Password: "secret",
	})

	if err != nil {
		t.Fatalf("LoginUser: %v", err)
	}
	if out.UserID != "user-1" {
		t.Errorf("UserID = %q, want user-1", out.UserID)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Error("tokens should be set")
	}
	if out.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}
}

func TestLoginUser_LoginPasswordRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := auth.NewAuthUseCase(userRepo)

	_, err := uc.LoginUser(context.Background(), auth.LoginUserInput{
		Login:    "",
		Password: "secret",
	})
	if !errors.Is(err, auth.ErrLoginPasswordRequired) {
		t.Errorf("err = %v, want ErrLoginPasswordRequired", err)
	}
}

func TestLoginUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "nobody").
		Return(nil, nil)

	uc := auth.NewAuthUseCase(userRepo)
	_, err := uc.LoginUser(context.Background(), auth.LoginUserInput{
		Login:    "nobody",
		Password: "secret",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Errorf("err = %v, want ErrInvalidCredentials", err)
	}
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, _ := crypto.HashPassword("correct")
	user := &models.User{ID: "user-1", Login: "testuser", PasswordHash: hash}

	userRepo := mocks.NewMockUserRepository(ctrl)
	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	uc := auth.NewAuthUseCase(userRepo)
	_, err := uc.LoginUser(context.Background(), auth.LoginUserInput{
		Login:    "testuser",
		Password: "wrong",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Errorf("err = %v, want ErrInvalidCredentials", err)
	}
}
