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

func TestGetData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	item := &models.Data{ID: "data-1", UserID: "user-1", Name: "test", Type: models.DataTypeText}
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Get(gomock.Any(), "user-1", "data-1").
		Return(item, nil)

	uc := data.NewDataUseCase(dataRepo)
	out, err := uc.GetData(context.Background(), data.GetDataInput{
		UserID: "user-1",
		DataID: "data-1",
	})

	if err != nil {
		t.Fatalf("GetData: %v", err)
	}
	if out.Data != item {
		t.Errorf("Data = %v, want %v", out.Data, item)
	}
	if out.Data.ID != "data-1" {
		t.Errorf("Data.ID = %q, want data-1", out.Data.ID)
	}
}

func TestGetData_DataIDRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	uc := data.NewDataUseCase(dataRepo)

	_, err := uc.GetData(context.Background(), data.GetDataInput{
		UserID: "user-1",
		DataID: "",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, data.ErrDataIDRequired) {
		t.Errorf("err = %v, want ErrDataIDRequired", err)
	}
}

func TestGetData_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Get(gomock.Any(), "user-1", "missing").
		Return(nil, nil)

	uc := data.NewDataUseCase(dataRepo)
	_, err := uc.GetData(context.Background(), data.GetDataInput{
		UserID: "user-1",
		DataID: "missing",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, data.ErrDataNotFound) {
		t.Errorf("err = %v, want ErrDataNotFound", err)
	}
}

func TestGetData_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := errors.New("db error")
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Get(gomock.Any(), "user-1", "data-1").
		Return(nil, wantErr)

	uc := data.NewDataUseCase(dataRepo)
	_, err := uc.GetData(context.Background(), data.GetDataInput{
		UserID: "user-1",
		DataID: "data-1",
	})

	if err != wantErr {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
