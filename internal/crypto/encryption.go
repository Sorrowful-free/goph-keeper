package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// SaltSize размер соли для PBKDF2
	SaltSize = 32
	// NonceSize размер nonce для AES-GCM
	NonceSize = 12
	// KeySize размер ключа для AES-256
	KeySize = 32
	// PBKDF2Iterations количество итераций PBKDF2
	PBKDF2Iterations = 100000
)

// EncryptData шифрует данные с использованием AES-256-GCM
// Использует PBKDF2 для получения ключа из пароля пользователя
func EncryptData(data []byte, password string) ([]byte, error) {
	// Генерируем соль
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// Генерируем ключ из пароля
	key := pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeySize, sha256.New)

	// Создаём AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Создаём GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Генерируем nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Шифруем данные
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	// Формат: salt + nonce + ciphertext
	result := make([]byte, 0, SaltSize+NonceSize+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptData расшифровывает данные
func DecryptData(encryptedData []byte, password string) ([]byte, error) {
	if len(encryptedData) < SaltSize+NonceSize {
		return nil, errors.New("encrypted data too short")
	}

	// Извлекаем соль, nonce и ciphertext
	salt := encryptedData[:SaltSize]
	nonce := encryptedData[SaltSize : SaltSize+NonceSize]
	ciphertext := encryptedData[SaltSize+NonceSize:]

	// Генерируем ключ из пароля
	key := pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeySize, sha256.New)

	// Создаём AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Создаём GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Расшифровываем данные
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
