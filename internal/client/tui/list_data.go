package tui

import (
	"fmt"
	"iter"
	"slices"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gophkeeper/gophkeeper/proto"
)

// ListDataModel представляет модель списка данных
type ListDataModel struct {
	model      *Model
	selected   int
	dataList   []*proto.Data
	err        error
	loading    bool
}

func NewListDataModel(m *Model) *ListDataModel {
	model := &ListDataModel{
		model:    m,
		selected: 0,
		loading:  true,
	}
	return model
}

func (m *ListDataModel) Init() tea.Cmd {
	return m.loadData()
}

func (m *ListDataModel) loadData() tea.Cmd {
	return func() tea.Msg {
		data, err := m.model.client.ListData(proto.DataType_UNKNOWN)
		if err != nil {
			return err
		}
		return data
	}
}

func (m *ListDataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []*proto.Data:
		if msg == nil {
			msg = []*proto.Data{}
		}
		m.dataList = msg
		m.loading = false
		m.model.dataList = msg
	case error:
		m.err = msg
		m.loading = false
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.dataList)-1 {
				m.selected++
			}
		case "enter":
			if len(m.dataList) > 0 && m.selected < len(m.dataList) {
				m.model.currentData = m.dataList[m.selected]
				m.model.state = StateViewData
				return NewViewDataModel(m.model), nil
			}
		case "r":
			m.loading = true
			return m, m.loadData()
		case "esc", "q":
			m.model.state = StateMainMenu
			return NewMainMenuModel(m.model), nil
		}
	}
	return m, nil
}

func (m *ListDataModel) View() string {
	if m.loading {
		return "Загрузка данных..."
	}

	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Ошибка: %v\n\nНажмите r для обновления, Esc для возврата", m.err))
	}

	if len(m.dataList) == 0 {
		return "Нет данных\n\nНажмите Esc для возврата"
	}

	items := []string{
		titleStyle.Render("Список данных"),
		"",
	}
	items = append(items, slices.Collect(listDataToLinesSeq(m.dataList, m.selected))...)
	items = append(items, "", "↑↓ для навигации, Enter для просмотра, r для обновления, Esc для возврата")
	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

// listDataToLinesSeq возвращает итератор строк списка данных (имя [тип], выбранный/обычный стиль)
func listDataToLinesSeq(dataList []*proto.Data, selected int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i, d := range dataList {
			item := fmt.Sprintf("%s [%s]", d.Name, getDataTypeName(d.Type))
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

func getDataTypeName(dt proto.DataType) string {
	switch dt {
	case proto.DataType_LOGIN_PASSWORD:
		return "Логин/Пароль"
	case proto.DataType_TEXT:
		return "Текст"
	case proto.DataType_BINARY:
		return "Бинарные"
	case proto.DataType_BANK_CARD:
		return "Банковская карта"
	default:
		return "Неизвестно"
	}
}
