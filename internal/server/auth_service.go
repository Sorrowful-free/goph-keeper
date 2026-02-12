package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophkeeper/gophkeeper/proto"
	"github.com/gophkeeper/gophkeeper/internal/usecase/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthService реализует gRPC-сервис аутентификации (delivery layer)
type AuthService struct {
	proto.UnimplementedAuthServiceServer
	authUC *auth.AuthUseCase
}

// NewAuthService создаёт новый сервис аутентификации
func NewAuthService(authUC *auth.AuthUseCase) *AuthService {
	return &AuthService{
		authUC: authUC,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req == nil {
		return &proto.RegisterResponse{Success: false, Message: "request is required"}, nil
	}

	out, err := s.authUC.RegisterUser(ctx, auth.RegisterUserInput{
		Login:    req.Login,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrLoginPasswordRequired):
			return &proto.RegisterResponse{
				Success: false,
				Message: "login and password are required",
			}, nil
		case errors.Is(err, auth.ErrUserAlreadyExists):
			return &proto.RegisterResponse{
				Success: false,
				Message: "user already exists",
			}, nil
		default:
			return &proto.RegisterResponse{
				Success: false,
				Message: fmt.Sprintf("error creating user: %v", err),
			}, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto.RegisterResponse{
		Success: true,
		Message: "user registered successfully",
		UserId:  out.UserID,
	}, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req == nil {
		return &proto.LoginResponse{Success: false, Message: "request is required"}, nil
	}

	out, err := s.authUC.LoginUser(ctx, auth.LoginUserInput{
		Login:    req.Login,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrLoginPasswordRequired):
			return &proto.LoginResponse{
				Success: false,
				Message: "login and password are required",
			}, nil
		case errors.Is(err, auth.ErrInvalidCredentials):
			return &proto.LoginResponse{
				Success: false,
				Message: "invalid login or password",
			}, nil
		default:
			return &proto.LoginResponse{
				Success: false,
				Message: "internal error",
			}, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto.LoginResponse{
		Success:      true,
		Message:      "login successful",
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
	}, nil
}

// RefreshToken обновляет access токен
func (s *AuthService) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	if req == nil || req.RefreshToken == "" {
		return &proto.RefreshTokenResponse{Success: false}, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	out, err := s.authUC.RefreshToken(ctx, auth.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrRefreshTokenRequired), errors.Is(err, auth.ErrInvalidRefreshToken):
			return &proto.RefreshTokenResponse{Success: false}, status.Error(codes.Unauthenticated, "invalid refresh token")
		default:
			return &proto.RefreshTokenResponse{Success: false}, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto.RefreshTokenResponse{
		Success:      true,
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
	}, nil
}
