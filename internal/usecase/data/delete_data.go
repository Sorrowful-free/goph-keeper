package data

import (
	"context"
	"errors"
)

var (
	ErrDataIDRequiredForDelete = errors.New("data_id is required")
)

// DeleteDataInput входные данные для удаления
type DeleteDataInput struct {
	UserID string
	DataID string
}

// DeleteData удаляет данные пользователя
func (uc *DataUseCase) DeleteData(ctx context.Context, in DeleteDataInput) error {
	if in.DataID == "" {
		return ErrDataIDRequiredForDelete
	}
	return uc.dataRepo.Delete(ctx, in.UserID, in.DataID)
}
