package storage

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Storage представляет хранилище данных
type Storage struct {
	db *gorm.DB
}

// NewStorage создаёт новое хранилище (миграции не выполняются — их нужно запускать отдельно через migrations.RunUp или CLI migrate).
func NewStorage(dsn string) (*Storage, error) {
	var db *gorm.DB
	var err error

	if dsn == "" {
		dsn = "gophkeeper.db"
	}

	if os.Getenv("DB_TYPE") == "postgres" || (len(dsn) >= 4 && dsn[:4] == "post") {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем, что соединение с БД действительно установлено
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

// GetDB возвращает экземпляр БД
func (s *Storage) GetDB() *gorm.DB {
	return s.db
}

// Close закрывает соединение с БД
func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
