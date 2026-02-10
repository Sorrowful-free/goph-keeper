package repository

import (
	"context"
	"time"

	domainrepo "github.com/gophkeeper/gophkeeper/internal/domain/repository"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/storage"
)

// dataRepo реализует domain/repository.DataRepository
type dataRepo struct {
	storage *storage.Storage
}

// NewDataRepository создаёт репозиторий данных
func NewDataRepository(storage *storage.Storage) domainrepo.DataRepository {
	return &dataRepo{storage: storage}
}

// Save сохраняет или обновляет данные
func (r *dataRepo) Save(ctx context.Context, userID string, data *models.Data) error {
	return r.storage.SaveData(userID, data)
}

// Get возвращает данные по ID
func (r *dataRepo) Get(ctx context.Context, userID, dataID string) (*models.Data, error) {
	return r.storage.GetData(userID, dataID)
}

// List возвращает список данных пользователя
func (r *dataRepo) List(ctx context.Context, userID string, dataType models.DataType) ([]*models.Data, error) {
	return r.storage.ListData(userID, dataType)
}

// Delete удаляет данные
func (r *dataRepo) Delete(ctx context.Context, userID, dataID string) error {
	return r.storage.DeleteData(userID, dataID)
}

// GetSince возвращает данные, изменённые после указанного времени
func (r *dataRepo) GetSince(ctx context.Context, userID string, since time.Time) ([]*models.Data, error) {
	return r.storage.GetDataSince(userID, since)
}
