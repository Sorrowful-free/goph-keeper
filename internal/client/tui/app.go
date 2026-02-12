package tui

import (
	"github.com/charmbracelet/bubbletea"
)

// AppModel представляет основную модель приложения
type AppModel struct {
	*Model
	current tea.Model
}

// NewAppModel создаёт новую модель приложения
func NewAppModel(serverAddress string) (*AppModel, error) {
	model, err := NewModel(serverAddress)
	if err != nil {
		return nil, err
	}

	app := &AppModel{
		Model:   model,
		current: NewLoginModel(model),
	}

	return app, nil
}

func (m *AppModel) Init() tea.Cmd {
	return m.current.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quit = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Обработка изменения размера окна
	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)

	// Проверяем, нужно ли переключить состояние
	if m.state == StateLogin && m.client.IsAuthenticated() {
		m.state = StateMainMenu
		m.current = NewMainMenuModel(m.Model)
		cmd = m.current.Init()
	}

	return m, cmd
}

func (m *AppModel) View() string {
	if m.quit {
		return ""
	}
	return m.current.View()
}
