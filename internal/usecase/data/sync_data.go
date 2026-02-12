package data

import (
	"context"
	"time"

	"github.com/gophkeeper/gophkeeper/internal/models"
)

// SyncDataInput входные данные для синхронизации
type SyncDataInput struct {
	UserID       string
	LastSyncTime time.Time
}

// SyncDataOutput результат синхронизации
type SyncDataOutput struct {
	Items    []*models.Data
	SyncTime time.Time
}

// SyncData возвращает данные, изменённые после указанного времени
func (uc *DataUseCase) SyncData(ctx context.Context, in SyncDataInput) (*SyncDataOutput, error) {
	var items []*models.Data
	var err error

	if in.LastSyncTime.IsZero() {
		items, err = uc.dataRepo.List(ctx, in.UserID, "")
	} else {
		items, err = uc.dataRepo.GetSince(ctx, in.UserID, in.LastSyncTime)
	}
	if err != nil {
		return nil, err
	}

	return &SyncDataOutput{
		Items:    items,
		SyncTime: time.Now(),
	}, nil
}
