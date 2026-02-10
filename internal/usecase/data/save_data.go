package data

import (
	"context"
	"errors"

	"github.com/gophkeeper/gophkeeper/internal/models"
)

var (
	ErrDataRequired = errors.New("data is required")
)

// SaveDataInput входные данные для сохранения
type SaveDataInput struct {
	UserID string
	Data   *models.Data
}

// SaveDataOutput результат сохранения
type SaveDataOutput struct {
	DataID  string
	Version int64
}

// SaveData сохраняет или обновляет данные пользователя
func (uc *DataUseCase) SaveData(ctx context.Context, in SaveDataInput) (*SaveDataOutput, error) {
	if in.Data == nil {
		return nil, ErrDataRequired
	}

	if err := uc.dataRepo.Save(ctx, in.UserID, in.Data); err != nil {
		return nil, err
	}

	return &SaveDataOutput{
		DataID:  in.Data.ID,
		Version: in.Data.Version,
	}, nil
}
