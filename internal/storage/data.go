package storage

import (
	"time"

	"github.com/gophkeeper/gophkeeper/internal/models"
	"gorm.io/gorm"
)

// SaveData сохраняет данные пользователя
func (s *Storage) SaveData(userID string, data *models.Data) error {
	// NOT NULL: пустой слайс вместо nil для encrypted_data
	if data.EncryptedData == nil {
		data.EncryptedData = []byte{}
	}

	// Проверяем существование записи
	var existingData models.Data
	err := s.db.Where("id = ? AND user_id = ?", data.ID, userID).First(&existingData).Error

	if err == nil {
		// Обновляем существующую запись
		data.Version = existingData.Version + 1
		data.UpdatedAt = time.Now()
		return s.db.Model(&existingData).Updates(data).Error
	} else if err == gorm.ErrRecordNotFound {
		// Создаём новую запись
		data.UserID = userID
		return s.db.Create(data).Error
	}

	return err
}

// GetData получает данные по ID
func (s *Storage) GetData(userID, dataID string) (*models.Data, error) {
	var data models.Data
	if err := s.db.Where("id = ? AND user_id = ?", dataID, userID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// ListData получает список данных пользователя
func (s *Storage) ListData(userID string, dataType models.DataType) ([]*models.Data, error) {
	var dataList []*models.Data
	query := s.db.Where("user_id = ?", userID)

	if dataType != "" {
		query = query.Where("type = ?", dataType)
	}

	if err := query.Find(&dataList).Error; err != nil {
		return nil, err
	}

	return dataList, nil
}

// DeleteData удаляет данные
func (s *Storage) DeleteData(userID, dataID string) error {
	return s.db.Where("id = ? AND user_id = ?", dataID, userID).Delete(&models.Data{}).Error
}

// GetDataSince получает данные, изменённые после указанного времени (для синхронизации)
func (s *Storage) GetDataSince(userID string, since time.Time) ([]*models.Data, error) {
	var dataList []*models.Data
	if err := s.db.Where("user_id = ? AND updated_at > ?", userID, since).Find(&dataList).Error; err != nil {
		return nil, err
	}
	return dataList, nil
}
