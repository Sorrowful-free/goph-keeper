package data

import (
	"context"
	"errors"

	"github.com/gophkeeper/gophkeeper/internal/models"
)

var (
	ErrDataIDRequired = errors.New("data_id is required")
	ErrDataNotFound   = errors.New("data not found")
)

// GetDataInput входные данные для получения записи
type GetDataInput struct {
	UserID string
	DataID string
}

// GetDataOutput результат получения данных
type GetDataOutput struct {
	Data *models.Data
}

// GetData возвращает данные по ID пользователя
func (uc *DataUseCase) GetData(ctx context.Context, in GetDataInput) (*GetDataOutput, error) {
	if in.DataID == "" {
		return nil, ErrDataIDRequired
	}

	data, err := uc.dataRepo.Get(ctx, in.UserID, in.DataID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrDataNotFound
	}

	return &GetDataOutput{Data: data}, nil
}
