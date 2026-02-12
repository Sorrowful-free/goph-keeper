package data_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gophkeeper/gophkeeper/internal/domain/repository/mocks"
	"github.com/gophkeeper/gophkeeper/internal/usecase/data"
	"go.uber.org/mock/gomock"
)

func TestDeleteData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Delete(gomock.Any(), "user-1", "data-1").
		Return(nil)

	uc := data.NewDataUseCase(dataRepo)
	err := uc.DeleteData(context.Background(), data.DeleteDataInput{
		UserID: "user-1",
		DataID: "data-1",
	})

	if err != nil {
		t.Fatalf("DeleteData: %v", err)
	}
}

func TestDeleteData_DataIDRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	uc := data.NewDataUseCase(dataRepo)

	err := uc.DeleteData(context.Background(), data.DeleteDataInput{
		UserID: "user-1",
		DataID: "",
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, data.ErrDataIDRequiredForDelete) {
		t.Errorf("err = %v, want ErrDataIDRequiredForDelete", err)
	}
}

func TestDeleteData_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := errors.New("db error")
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		Delete(gomock.Any(), "user-1", "data-1").
		Return(wantErr)

	uc := data.NewDataUseCase(dataRepo)
	err := uc.DeleteData(context.Background(), data.DeleteDataInput{
		UserID: "user-1",
		DataID: "data-1",
	})

	if err != wantErr {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
