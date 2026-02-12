package repository

import (
	"context"
	"time"

	"github.com/gophkeeper/gophkeeper/internal/models"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_data_repository.go -package=mocks github.com/gophkeeper/gophkeeper/internal/domain/repository DataRepository

// DataRepository определяет контракт для работы с данными пользователя
type DataRepository interface {
	Save(ctx context.Context, userID string, data *models.Data) error
	Get(ctx context.Context, userID, dataID string) (*models.Data, error)
	List(ctx context.Context, userID string, dataType models.DataType) ([]*models.Data, error)
	Delete(ctx context.Context, userID, dataID string) error
	GetSince(ctx context.Context, userID string, since time.Time) ([]*models.Data, error)
}
