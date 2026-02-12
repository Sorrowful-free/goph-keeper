package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1).
			Width(40)

	focusedStyle = inputStyle.Copy().
			BorderForeground(lipgloss.Color("205"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			MarginTop(1)
)

// Режим экрана: вход или регистрация
const (
	loginMode = iota
	registerMode
)

// LoginModel представляет модель входа и регистрации
type LoginModel struct {
	model         *Model
	loginInput    textinput.Model
	passwordInput textinput.Model
	focused       int
	mode          int // loginMode или registerMode
	err           error
}

func NewLoginModel(m *Model) *LoginModel {
	loginInput := textinput.New()
	loginInput.Placeholder = "Логин"
	loginInput.Focus()
	loginInput.CharLimit = 50
	loginInput.Width = 38

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Пароль"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.CharLimit = 50
	passwordInput.Width = 38

	return &LoginModel{
		model:         m,
		loginInput:    loginInput,
		passwordInput: passwordInput,
		focused:       0,
		mode:          loginMode,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "enter", "up", "down":
			if msg.String() == "enter" {
				return m.handleSubmit()
			}
			if msg.String() == "tab" || msg.String() == "down" {
				m.focused = (m.focused + 1) % 2
			} else if msg.String() == "up" {
				m.focused = (m.focused + 2 - 1) % 2
			}
			if m.focused == 0 {
				m.loginInput.Focus()
				m.passwordInput.Blur()
			} else {
				m.loginInput.Blur()
				m.passwordInput.Focus()
			}
		case "left", "right", "1", "2":
			// Переключение сценария: Вход (1) / Регистрация (2) или стрелки
			if msg.String() == "right" || msg.String() == "2" {
				m.mode = registerMode
			} else {
				m.mode = loginMode
			}
			m.err = nil
		case "ctrl+c", "q":
			m.model.quit = true
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	if m.focused == 0 {
		m.loginInput, cmd = m.loginInput.Update(msg)
	} else {
		m.passwordInput, cmd = m.passwordInput.Update(msg)
	}

	return m, cmd
}

func (m *LoginModel) handleSubmit() (tea.Model, tea.Cmd) {
	login := m.loginInput.Value()
	password := m.passwordInput.Value()

	if login == "" || password == "" {
		m.err = fmt.Errorf("логин и пароль обязательны")
		return m, nil
	}

	if m.mode == registerMode {
		// Сценарий регистрации
		if err := m.model.client.Register(login, password); err != nil {
			m.err = err
			return m, nil
		}
		// После успешной регистрации автоматически входим
		if err := m.model.client.Login(login, password); err != nil {
			m.err = err
			return m, nil
		}
	} else {
		// Сценарий входа
		if err := m.model.client.Login(login, password); err != nil {
			m.err = err
			return m, nil
		}
	}

	// Успешный вход — переходим в главное меню
	m.model.state = StateMainMenu
	return NewMainMenuModel(m.model), nil
}

func (m *LoginModel) View() string {
	var style lipgloss.Style
	if m.focused == 0 {
		style = focusedStyle
	} else {
		style = inputStyle
	}

	loginView := style.Render(m.loginInput.View())

	if m.focused == 1 {
		style = focusedStyle
	} else {
		style = inputStyle
	}

	passwordView := style.Render(m.passwordInput.View())

	title := "Вход в GophKeeper"
	if m.mode == registerMode {
		title = "Регистрация в GophKeeper"
	}
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		title,
		loginView,
		passwordView,
	)

	if m.err != nil {
		view += errorStyle.Render(fmt.Sprintf("Ошибка: %v", m.err))
	}

	modeHint := "← Вход | Регистрация →"
	if m.mode == registerMode {
		modeHint = "Вход ← | Регистрация →"
	}
	view += fmt.Sprintf("\n\nРежим: %s", modeHint)
	view += "\nTab — поля, Enter — отправить, ←/→ или 1/2 — смена режима, q — выход"

	return lipgloss.NewStyle().Padding(1, 2).Render(view)
}
