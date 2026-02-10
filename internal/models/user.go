package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User представляет пользователя системы
type User struct {
	ID           string         `gorm:"primaryKey;size:36" json:"id"`
	Login        string         `gorm:"uniqueIndex;not null" json:"login"`
	PasswordHash string         `gorm:"not null" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate генерирует UUID для новых пользователей (совместимо с SQLite и PostgreSQL)
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// TableName возвращает имя таблицы
func (User) TableName() string {
	return "users"
}
