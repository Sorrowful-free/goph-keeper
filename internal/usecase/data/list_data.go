package data

import (
	"context"

	"github.com/gophkeeper/gophkeeper/internal/models"
)

// ListDataInput входные данные для списка
type ListDataInput struct {
	UserID   string
	DataType models.DataType
}

// ListDataOutput результат списка данных
type ListDataOutput struct {
	Items []*models.Data
}

// ListData возвращает список данных пользователя (опционально по типу)
func (uc *DataUseCase) ListData(ctx context.Context, in ListDataInput) (*ListDataOutput, error) {
	items, err := uc.dataRepo.List(ctx, in.UserID, in.DataType)
	if err != nil {
		return nil, err
	}
	return &ListDataOutput{Items: items}, nil
}
