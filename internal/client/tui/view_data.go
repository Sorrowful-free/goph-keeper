package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gophkeeper/gophkeeper/internal/client/format"
)

// ViewDataModel представляет модель просмотра данных
type ViewDataModel struct {
	model *Model
}

func NewViewDataModel(m *Model) *ViewDataModel {
	return &ViewDataModel{model: m}
}

func (m *ViewDataModel) Init() tea.Cmd {
	return nil
}

func (m *ViewDataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.model.state = StateListData
			listModel := NewListDataModel(m.model)
			return listModel, listModel.Init()
		case "d":
			// Удаление данных
			if m.model.currentData != nil {
				m.model.state = StateDeleteData
				return NewDeleteDataModel(m.model), nil
			}
		}
	}
	return m, nil
}

func (m *ViewDataModel) View() string {
	if m.model.currentData == nil {
		return "Нет данных для отображения"
	}

	data := m.model.currentData
	var view []string

	view = append(view, titleStyle.Render("Просмотр данных"))
	view = append(view, "")
	view = append(view, fmt.Sprintf("Название: %s", data.Name))
	view = append(view, fmt.Sprintf("Тип: %s", format.DataTypeDisplayName(data.Type)))
	view = append(view, fmt.Sprintf("ID: %s", data.Id))
	view = append(view, fmt.Sprintf("Версия: %d", data.Version))
	view = append(view, "")

	if len(data.Metadata) > 0 {
		view = append(view, "Метаданные:")
		view = append(view, format.MetadataToDisplayLines(data.Metadata)...)
		view = append(view, "")
	}

	if content := format.DataContentToDisplayString(data); content != "" {
		view = append(view, "Данные:")
		view = append(view, content)
		view = append(view, "")
	} else if len(data.EncryptedData) > 0 {
		view = append(view, fmt.Sprintf("Данные (сырые): %d байт", len(data.EncryptedData)))
		view = append(view, "")
	}

	view = append(view, "Esc для возврата, d для удаления")

	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, view...))
}
