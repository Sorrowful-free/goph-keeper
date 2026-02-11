package tui

import (
	"iter"
	"slices"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Margin(1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			MarginBottom(1)

	selectedMenuItemStyle = menuItemStyle.Copy().
				Foreground(lipgloss.Color("205")).
				Bold(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			MarginBottom(1)
)

// MainMenuModel представляет модель главного меню
type MainMenuModel struct {
	model      *Model
	selected   int
	menuItems  []string
}

func NewMainMenuModel(m *Model) *MainMenuModel {
	return &MainMenuModel{
		model:     m,
		selected:  0,
		menuItems: []string{
			"📋 Список данных",
			"➕ Добавить данные",
			"🔄 Синхронизация",
			"🚪 Выход",
		},
	}
}

func (m *MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m *MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.menuItems)-1 {
				m.selected++
			}
		case "enter":
			return m.handleSelection()
		case "q", "ctrl+c":
			m.model.quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *MainMenuModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.selected {
	case 0: // Список данных
		m.model.state = StateListData
		listModel := NewListDataModel(m.model)
		return listModel, listModel.Init()
	case 1: // Добавить данные
		m.model.state = StateAddData
		return NewAddDataModel(m.model), nil
	case 2: // Синхронизация
		m.model.state = StateSync
		syncModel := NewSyncModel(m.model)
		return syncModel, syncModel.Init()
	case 3: // Выход
		m.model.quit = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *MainMenuModel) View() string {
	items := []string{
		titleStyle.Render("GophKeeper - Главное меню"),
		"",
	}
	items = append(items, slices.Collect(menuItemsToLinesSeq(m.menuItems, m.selected))...)
	items = append(items, "", "Используйте ↑↓ для навигации, Enter для выбора, q для выхода")
	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

// menuItemsToLinesSeq возвращает итератор строк меню (выбранный/обычный стиль)
func menuItemsToLinesSeq(menuItems []string, selected int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i, item := range menuItems {
			var line string
			if i == selected {
				line = selectedMenuItemStyle.Render("▶ " + item)
			} else {
				line = menuItemStyle.Render("  " + item)
			}
			if !yield(line) {
				return
			}
		}
	}
}
