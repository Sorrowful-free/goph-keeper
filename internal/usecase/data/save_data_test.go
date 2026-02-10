package data_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gophkeeper/gophkeeper/internal/domain/repository/mocks"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/usecase/data"
	"go.uber.org/mock/gomock"
)

func TestSaveData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Save(gomock.Any(), "user-1", gomock.Any()).
		Return(nil)

	uc := data.NewDataUseCase(dataRepo)
	item := &models.Data{ID: "data-1", Name: "test", Type: models.DataTypeText, Version: 1}
	out, err := uc.SaveData(context.Background(), data.SaveDataInput{
		UserID: "user-1",
		Data:   item,
	})

	if err != nil {
		t.Fatalf("SaveData: %v", err)
	}
	if out.DataID != "data-1" {
		t.Errorf("DataID = %q, want data-1", out.DataID)
	}
	if out.Version != 1 {
		t.Errorf("Version = %d, want 1", out.Version)
	}
}

func TestSaveData_DataRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	uc := data.NewDataUseCase(dataRepo)

	_, err := uc.SaveData(context.Background(), data.SaveDataInput{
		UserID: "user-1",
		Data:   nil,
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, data.ErrDataRequired) {
		t.Errorf("err = %v, want ErrDataRequired", err)
	}
}

func TestSaveData_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := errors.New("db error")
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Save(gomock.Any(), "user-1", gomock.Any()).
		Return(wantErr)

	uc := data.NewDataUseCase(dataRepo)
	_, err := uc.SaveData(context.Background(), data.SaveDataInput{
		UserID: "user-1",
		Data:   &models.Data{ID: "data-1", Name: "test", Type: models.DataTypeText},
	})

	if err != wantErr {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
