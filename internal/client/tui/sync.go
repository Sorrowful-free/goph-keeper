package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gophkeeper/gophkeeper/proto"
)

// SyncModel представляет модель синхронизации
type SyncModel struct {
	model   *Model
	loading bool
	err     error
	message string
}

func NewSyncModel(m *Model) *SyncModel {
	return &SyncModel{
		model:   m,
		loading: true,
	}
}

func (m *SyncModel) Init() tea.Cmd {
	return m.sync()
}

func (m *SyncModel) sync() tea.Cmd {
	return func() tea.Msg {
		// Синхронизируем данные (lastSyncTime = 0 означает получить все данные)
		data, syncTime, err := m.model.client.SyncData(0)
		if err != nil {
			return err
		}
		return syncResult{data: data, syncTime: syncTime}
	}
}

type syncResult struct {
	data     []*proto.Data
	syncTime int64
}

func (m *SyncModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case syncResult:
		data := msg.data
		if data == nil {
			data = []*proto.Data{}
		}
		m.model.dataList = data
		m.loading = false
		if len(data) == 0 {
			m.message = "Данных нет. Синхронизация завершена."
		} else {
			m.message = fmt.Sprintf("Синхронизировано %d записей", len(data))
		}
	case error:
		m.err = msg
		m.loading = false
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			m.err = nil
			m.message = ""
			return m, m.sync()
		case "esc", "q":
			m.model.state = StateMainMenu
			return NewMainMenuModel(m.model), nil
		}
	}
	return m, nil
}

func (m *SyncModel) View() string {
	if m.loading {
		return "Синхронизация данных..."
	}

	var view []string
	view = append(view, titleStyle.Render("Синхронизация"))
	view = append(view, "")

	if m.err != nil {
		view = append(view, errorStyle.Render(fmt.Sprintf("Ошибка: %v", m.err)))
	} else if m.message != "" {
		view = append(view, successStyle.Render(m.message))
	} else {
		view = append(view, "Синхронизация завершена")
	}

	view = append(view, "")
	view = append(view, "r для повторной синхронизации, Esc для возврата")

	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, view...))
}
