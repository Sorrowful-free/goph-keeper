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

func TestListData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	items := []*models.Data{
		{ID: "data-1", UserID: "user-1", Name: "a", Type: models.DataTypeText},
		{ID: "data-2", UserID: "user-1", Name: "b", Type: models.DataTypeText},
	}
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		List(gomock.Any(), "user-1", models.DataTypeText).
		Return(items, nil)

	uc := data.NewDataUseCase(dataRepo)
	out, err := uc.ListData(context.Background(), data.ListDataInput{
		UserID:   "user-1",
		DataType: models.DataTypeText,
	})

	if err != nil {
		t.Fatalf("ListData: %v", err)
	}
	if len(out.Items) != 2 {
		t.Errorf("len(Items) = %d, want 2", len(out.Items))
	}
	if out.Items[0].ID != "data-1" || out.Items[1].ID != "data-2" {
		t.Errorf("Items IDs = %q, %q", out.Items[0].ID, out.Items[1].ID)
	}
}

func TestListData_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		List(gomock.Any(), "user-1", models.DataType("")).
		Return(nil, nil)

	uc := data.NewDataUseCase(dataRepo)
	out, err := uc.ListData(context.Background(), data.ListDataInput{
		UserID:   "user-1",
		DataType: "",
	})

	if err != nil {
		t.Fatalf("ListData: %v", err)
	}
	if out.Items != nil {
		t.Errorf("Items = %v, want nil", out.Items)
	}
}

func TestListData_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := errors.New("db error")
	dataRepo := mocks.NewMockDataRepository(ctrl)
	dataRepo.EXPECT().
		List(gomock.Any(), "user-1", gomock.Any()).
		Return(nil, wantErr)

	uc := data.NewDataUseCase(dataRepo)
	_, err := uc.ListData(context.Background(), data.ListDataInput{
		UserID:   "user-1",
		DataType: models.DataTypeText,
	})

	if err != wantErr {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}
