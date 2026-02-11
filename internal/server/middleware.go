package server

import (
	"context"
	"log"
	"time"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

// LoggingInterceptor логирует каждый входящий gRPC-запрос: метод, длительность и результат.
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Printf("grpc method=%s duration=%s code=%s msg=%s",
				info.FullMethod, duration, st.Code(), st.Message())
		} else {
			log.Printf("grpc method=%s duration=%s error=%v", info.FullMethod, duration, err)
		}
		return resp, err
	}

	log.Printf("grpc method=%s duration=%s code=OK", info.FullMethod, duration)
	return resp, nil
}

// AuthInterceptor перехватывает запросы, проверяет JWT и кладёт userID в контекст.
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Пропускаем аутентификацию для методов AuthService
	if info.FullMethod == "/gophkeeper.AuthService/Register" ||
		info.FullMethod == "/gophkeeper.AuthService/Login" ||
		info.FullMethod == "/gophkeeper.AuthService/RefreshToken" {
		return handler(ctx, req)
	}

	// Для остальных методов проверяем токен и извлекаем userID
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no authorization token")
	}

	token := tokens[0]
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := crypto.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	ctx = context.WithValue(ctx, userIDContextKey, claims.UserID)
	return handler(ctx, req)
}

// GetUserIDFromContext возвращает user ID, записанный в контекст AuthInterceptor.
// Используется в обработчиках (например, DataService) для получения текущего пользователя.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	v := ctx.Value(userIDContextKey)
	if v == nil {
		return "", status.Error(codes.Unauthenticated, "no user in context")
	}
	userID, ok := v.(string)
	if !ok {
		return "", status.Error(codes.Internal, "invalid user id in context")
	}
	return userID, nil
}
