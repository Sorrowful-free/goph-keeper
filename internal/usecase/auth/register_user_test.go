package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gophkeeper/gophkeeper/internal/domain/repository/mocks"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/usecase/auth"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(nil, nil)
	userRepo.EXPECT().
		Create(gomock.Any(), "testuser", gomock.Any()).
		Return(&models.User{ID: "user-1", Login: "testuser", PasswordHash: "hash"}, nil)

	uc := auth.NewAuthUseCase(userRepo)
	out, err := uc.RegisterUser(context.Background(), auth.RegisterUserInput{
		Login:    "testuser",
		Password: "secret",
	})

	if err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}
	if out.UserID != "user-1" {
		t.Errorf("UserID = %q, want user-1", out.UserID)
	}
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(&models.User{ID: "existing", Login: "testuser"}, nil)
	// Create не должен вызываться

	uc := auth.NewAuthUseCase(userRepo)
	_, err := uc.RegisterUser(context.Background(), auth.RegisterUserInput{
		Login:    "testuser",
		Password: "secret",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, auth.ErrUserAlreadyExists) {
		t.Errorf("err = %v, want ErrUserAlreadyExists", err)
	}
}

func TestRegisterUser_LoginPasswordRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	// Репозиторий не вызывается

	uc := auth.NewAuthUseCase(userRepo)

	t.Run("empty_login", func(t *testing.T) {
		_, err := uc.RegisterUser(context.Background(), auth.RegisterUserInput{
			Login:    "",
			Password: "secret",
		})
		if !errors.Is(err, auth.ErrLoginPasswordRequired) {
			t.Errorf("err = %v, want ErrLoginPasswordRequired", err)
		}
	})
	t.Run("empty_password", func(t *testing.T) {
		_, err := uc.RegisterUser(context.Background(), auth.RegisterUserInput{
			Login:    "user",
			Password: "",
		})
		if !errors.Is(err, auth.ErrLoginPasswordRequired) {
			t.Errorf("err = %v, want ErrLoginPasswordRequired", err)
		}
	})
}
