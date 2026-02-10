package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost определяет сложность хеширования паролей
	BcryptCost = 12
)

// HashPassword создаёт bcrypt хеш пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword проверяет соответствие пароля хешу
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
