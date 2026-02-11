package crypto

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Глобальные параметры JWT — задаются из конфига при старте сервера (см. cmd/server/main.go).
var (
	JWTSecret          []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
)

// Claims представляет JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func checkJWTConfig() error {
	if len(JWTSecret) == 0 {
		return errors.New("JWT secret not configured (set from config at server startup)")
	}
	return nil
}

// GenerateAccessToken генерирует access токен
func GenerateAccessToken(userID string) (string, error) {
	if err := checkJWTConfig(); err != nil {
		return "", err
	}
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// GenerateRefreshToken генерирует refresh токен
func GenerateRefreshToken(userID string) (string, error) {
	if err := checkJWTConfig(); err != nil {
		return "", err
	}
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken валидирует JWT токен
func ValidateToken(tokenString string) (*Claims, error) {
	if err := checkJWTConfig(); err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
