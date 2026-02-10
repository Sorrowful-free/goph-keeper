package data_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gophkeeper/gophkeeper/internal/domain/repository/mocks"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/usecase/data"
	"go.uber.org/mock/gomock"
)

func TestSyncData_AllData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	items := []*models.Data{
		{ID: "data-1", UserID: "user-1", Name: "a", Type: models.DataTypeText},
	}
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		List(gomock.Any(), "user-1", models.DataType("")).
		Return(items, nil)

	uc := data.NewDataUseCase(dataRepo)
	out, err := uc.SyncData(context.Background(), data.SyncDataInput{
		UserID:       "user-1",
		LastSyncTime: time.Time{},
	})

	if err != nil {
		t.Fatalf("SyncData: %v", err)
	}
	if len(out.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(out.Items))
	}
	if out.Items[0].ID != "data-1" {
		t.Errorf("Items[0].ID = %q, want data-1", out.Items[0].ID)
	}
	if out.SyncTime.IsZero() {
		t.Error("SyncTime should be set")
	}
}

func TestSyncData_Since(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	since := time.Now().Add(-time.Hour)
	items := []*models.Data{
		{ID: "data-2", UserID: "user-1", Name: "b", Type: models.DataTypeText},
	}
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		GetSince(gomock.Any(), "user-1", since).
		Return(items, nil)

	uc := data.NewDataUseCase(dataRepo)
	out, err := uc.SyncData(context.Background(), data.SyncDataInput{
		UserID:       "user-1",
		LastSyncTime: since,
	})

	if err != nil {
		t.Fatalf("SyncData: %v", err)
	}
	if len(out.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(out.Items))
	}
	if out.Items[0].ID != "data-2" {
		t.Errorf("Items[0].ID = %q, want data-2", out.Items[0].ID)
	}
}

func TestSyncData_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := errors.New("db error")
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		List(gomock.Any(), "user-1", models.DataType("")).
		Return(nil, wantErr)

	uc := data.NewDataUseCase(dataRepo)
	_, err := uc.SyncData(context.Background(), data.SyncDataInput{
		UserID:       "user-1",
		LastSyncTime: time.Time{},
	})

	if err != wantErr {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
