package storage

import (
	"errors"

	"github.com/gophkeeper/gophkeeper/internal/models"
	"gorm.io/gorm"
)

// CreateUser создаёт нового пользователя
func (s *Storage) CreateUser(login, passwordHash string) (*models.User, error) {
	user := &models.User{
		Login:        login,
		PasswordHash: passwordHash,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByLogin получает пользователя по логину
func (s *Storage) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (s *Storage) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
