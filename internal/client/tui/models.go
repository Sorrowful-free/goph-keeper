package tui

import (
	"github.com/gophkeeper/gophkeeper/internal/client"
	"github.com/gophkeeper/gophkeeper/proto"
)

// AppState представляет состояние приложения
type AppState int

const (
	StateLogin AppState = iota
	StateRegister
	StateMainMenu
	StateListData
	StateViewData
	StateAddData
	StateEditData
	StateDeleteData
	StateSync
	StateQuit
)

// Model представляет основную модель TUI
type Model struct {
	client      *client.Client
	state       AppState
	err         error
	message     string
	login       string
	password    string
	selectedIdx int
	dataList    []*proto.Data
	currentData *proto.Data
	quit        bool
}

// NewModel создаёт новую модель
func NewModel(serverAddress string) (*Model, error) {
	c, err := client.NewClient(serverAddress)
	if err != nil {
		return nil, err
	}

	return &Model{
		client:      c,
		state:       StateLogin,
		selectedIdx: 0,
	}, nil
}

// Close закрывает клиент
func (m *Model) Close() error {
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}
