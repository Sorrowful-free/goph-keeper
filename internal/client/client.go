package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gophkeeper/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client представляет клиент для взаимодействия с сервером
type Client struct {
	conn          *grpc.ClientConn
	authClient    proto.AuthServiceClient
	dataClient    proto.DataServiceClient
	accessToken   string
	refreshToken  string
	serverAddress string
}

// NewClient создаёт новый клиент
func NewClient(serverAddress string) (*Client, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return &Client{
		conn:          conn,
		authClient:    proto.NewAuthServiceClient(conn),
		dataClient:    proto.NewDataServiceClient(conn),
		serverAddress: serverAddress,
	}, nil
}

// NewClientWithClients создаёт клиент с заданными gRPC-клиентами (для тестов).
func NewClientWithClients(authClient proto.AuthServiceClient, dataClient proto.DataServiceClient) *Client {
	return &Client{
		authClient: authClient,
		dataClient: dataClient,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Register регистрирует нового пользователя
func (c *Client) Register(login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Register(ctx, &proto.RegisterRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("registration failed: %s", resp.Message)
	}

	return nil
}

// Login выполняет вход
func (c *Client) Login(login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Login(ctx, &proto.LoginRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("login failed: %s", resp.Message)
	}

	c.accessToken = resp.AccessToken
	c.refreshToken = resp.RefreshToken

	return nil
}

// RefreshToken обновляет токен
func (c *Client) RefreshToken() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.RefreshToken(ctx, &proto.RefreshTokenRequest{
		RefreshToken: c.refreshToken,
	})

	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("token refresh failed")
	}

	c.accessToken = resp.AccessToken
	c.refreshToken = resp.RefreshToken

	return nil
}

// getContext создаёт контекст с токеном авторизации
func (c *Client) getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + c.accessToken,
	})
	return metadata.NewOutgoingContext(ctx, md), cancel
}

// IsAuthenticated проверяет, аутентифицирован ли клиент
func (c *Client) IsAuthenticated() bool {
	return c.accessToken != ""
}

// SaveData сохраняет данные
func (c *Client) SaveData(data *proto.Data) (string, int64, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	resp, err := c.dataClient.SaveData(ctx, &proto.SaveDataRequest{
		Data: data,
	})

	if err != nil {
		return "", 0, err
	}

	if !resp.Success {
		return "", 0, fmt.Errorf("save failed: %s", resp.Message)
	}

	return resp.DataId, resp.Version, nil
}

// GetData получает данные по ID
func (c *Client) GetData(dataID string) (*proto.Data, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	resp, err := c.dataClient.GetData(ctx, &proto.GetDataRequest{
		DataId: dataID,
	})

	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("get failed: %s", resp.Message)
	}

	return resp.Data, nil
}

// ListData получает список данных
func (c *Client) ListData(dataType proto.DataType) ([]*proto.Data, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	resp, err := c.dataClient.ListData(ctx, &proto.ListDataRequest{
		Type: dataType,
	})

	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("list failed: %s", resp.Message)
	}

	if resp.Data == nil {
		return []*proto.Data{}, nil
	}
	return resp.Data, nil
}

// DeleteData удаляет данные
func (c *Client) DeleteData(dataID string) error {
	ctx, cancel := c.getContext()
	defer cancel()

	resp, err := c.dataClient.DeleteData(ctx, &proto.DeleteDataRequest{
		DataId: dataID,
	})

	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("delete failed: %s", resp.Message)
	}

	return nil
}

// SyncData синхронизирует данные
func (c *Client) SyncData(lastSyncTime int64) ([]*proto.Data, int64, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	resp, err := c.dataClient.SyncData(ctx, &proto.SyncDataRequest{
		LastSyncTime: lastSyncTime,
	})

	if err != nil {
		return nil, 0, err
	}

	if !resp.Success {
		return nil, 0, fmt.Errorf("sync failed: %s", resp.Message)
	}

	data := resp.Data
	if data == nil {
		data = []*proto.Data{}
	}
	return data, resp.SyncTime, nil
}
