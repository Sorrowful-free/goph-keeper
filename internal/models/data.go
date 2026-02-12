package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DataType представляет тип хранимых данных
type DataType string

const (
	DataTypeLoginPassword DataType = "login_password"
	DataTypeText          DataType = "text"
	DataTypeBinary        DataType = "binary"
	DataTypeBankCard      DataType = "bank_card"
)

// Data представляет хранимые данные пользователя
type Data struct {
	ID            string         `gorm:"primaryKey;size:36" json:"id"`
	UserID        string         `gorm:"size:36;not null;index" json:"user_id"`
	Type          DataType       `gorm:"size:50;not null;index" json:"type"`
	Name          string         `gorm:"not null" json:"name"`
	EncryptedData []byte         `gorm:"not null" json:"-"` // blob в SQLite, bytea в PostgreSQL
	Metadata      string         `gorm:"type:text" json:"metadata"` // JSON строка для метаданных
	Version       int64          `gorm:"default:1" json:"version"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate генерирует UUID для новых записей (совместимо с SQLite и PostgreSQL)
func (d *Data) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// TableName возвращает имя таблицы
func (Data) TableName() string {
	return "data"
}

// MetadataItem представляет элемент метаданных
type MetadataItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
