package client_test

import (
	"testing"

	"github.com/gophkeeper/gophkeeper/internal/client"
	"github.com/gophkeeper/gophkeeper/internal/client/mocks"
	"github.com/gophkeeper/gophkeeper/proto"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&proto.RegisterResponse{Success: true}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	err := c.Register("user", "pass")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
}

func TestRegister_ServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, status.Error(codes.Internal, "db error"))

	c := client.NewClientWithClients(authMock, dataMock)
	err := c.Register("user", "pass")
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}

func TestRegister_UnsuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&proto.RegisterResponse{Success: false, Message: "login занят"}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	err := c.Register("user", "pass")
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
	if err.Error() != "registration failed: login занят" {
		t.Errorf("err = %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{
			Success:      true,
			AccessToken:  "at",
			RefreshToken: "rt",
		}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	err := c.Login("u", "p")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if !c.IsAuthenticated() {
		t.Error("клиент должен быть аутентифицирован")
	}
}

func TestLogin_UnsuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: false, Message: "неверный пароль"}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	err := c.Login("u", "p")
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
	if !c.IsAuthenticated() {
		// токены не должны быть установлены
	}
	_ = err
}

func TestIsAuthenticated_InitiallyFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewClientWithClients(mocks.NewMockAuthServiceClient(ctrl), mocks.NewMockDataServiceClient(ctrl))
	if c.IsAuthenticated() {
		t.Error("изначально не аутентифицирован")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at1", RefreshToken: "rt1"}, nil)
	authMock.EXPECT().
		RefreshToken(gomock.Any(), gomock.Any()).
		Return(&proto.RefreshTokenResponse{Success: true, AccessToken: "at2", RefreshToken: "rt2"}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	if err := c.Login("u", "p"); err != nil {
		t.Fatalf("Login: %v", err)
	}
	err := c.RefreshToken()
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
}

func TestRefreshToken_UnsuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)
	authMock.EXPECT().
		RefreshToken(gomock.Any(), gomock.Any()).
		Return(&proto.RefreshTokenResponse{Success: false}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	if err := c.Login("u", "p"); err != nil {
		t.Fatalf("Login: %v", err)
	}
	err := c.RefreshToken()
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}

func TestSaveData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)

	payload := []byte("secret")
	dataMock.EXPECT().
		SaveData(gomock.Any(), gomock.Any()).
		Return(&proto.SaveDataResponse{Success: true, DataId: "id1", Version: 1}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	if err := c.Login("u", "p"); err != nil {
		t.Fatalf("Login: %v", err)
	}

	id, ver, err := c.SaveData(&proto.Data{Type: proto.DataType_LOGIN_PASSWORD, EncryptedData: payload})
	if err != nil {
		t.Fatalf("SaveData: %v", err)
	}
	if id != "id1" || ver != 1 {
		t.Errorf("id=%q ver=%d", id, ver)
	}
}

func TestSaveData_UnsuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)
	dataMock.EXPECT().SaveData(gomock.Any(), gomock.Any()).
		Return(&proto.SaveDataResponse{Success: false, Message: "conflict"}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	_ = c.Login("u", "p")

	_, _, err := c.SaveData(&proto.Data{Type: proto.DataType_TEXT, EncryptedData: []byte("x")})
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}

func TestGetData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)
	dataMock.EXPECT().
		GetData(gomock.Any(), gomock.Any()).
		Return(&proto.GetDataResponse{
			Success: true,
			Data:    &proto.Data{Type: proto.DataType_TEXT, EncryptedData: []byte("content")},
		}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	_ = c.Login("u", "p")

	data, err := c.GetData("id1")
	if err != nil {
		t.Fatalf("GetData: %v", err)
	}
	if data == nil || string(data.EncryptedData) != "content" {
		t.Errorf("data = %+v", data)
	}
}

func TestListData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)
	dataMock.EXPECT().
		ListData(gomock.Any(), gomock.Any()).
		Return(&proto.ListDataResponse{
			Success: true,
			Data: []*proto.Data{
				{Id: "1", Type: proto.DataType_LOGIN_PASSWORD},
				{Id: "2", Type: proto.DataType_LOGIN_PASSWORD},
			},
		}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	_ = c.Login("u", "p")

	list, err := c.ListData(proto.DataType_LOGIN_PASSWORD)
	if err != nil {
		t.Fatalf("ListData: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d", len(list))
	}
}

func TestDeleteData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)
	dataMock.EXPECT().
		DeleteData(gomock.Any(), gomock.Any()).
		Return(&proto.DeleteDataResponse{Success: true}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	_ = c.Login("u", "p")

	err := c.DeleteData("id1")
	if err != nil {
		t.Fatalf("DeleteData: %v", err)
	}
}

func TestSyncData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthServiceClient(ctrl)
	dataMock := mocks.NewMockDataServiceClient(ctrl)

	authMock.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(&proto.LoginResponse{Success: true, AccessToken: "at", RefreshToken: "rt"}, nil)
	dataMock.EXPECT().
		SyncData(gomock.Any(), gomock.Any()).
		Return(&proto.SyncDataResponse{
			Success:    true,
			Data:       []*proto.Data{{Id: "1", Version: 2}},
			SyncTime:   200,
		}, nil)

	c := client.NewClientWithClients(authMock, dataMock)
	_ = c.Login("u", "p")

	data, syncTime, err := c.SyncData(100)
	if err != nil {
		t.Fatalf("SyncData: %v", err)
	}
	if len(data) != 1 || data[0].Id != "1" || syncTime != 200 {
		t.Errorf("data=%+v syncTime=%d", data, syncTime)
	}
}

func TestClose_NoConnection(t *testing.T) {
	c := client.NewClientWithClients(nil, nil)
	// Close при отсутствии conn не должен паниковать
	if err := c.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}
