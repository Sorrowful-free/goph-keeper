package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DeleteDataModel представляет модель удаления данных
type DeleteDataModel struct {
	model *Model
}

func NewDeleteDataModel(m *Model) *DeleteDataModel {
	return &DeleteDataModel{model: m}
}

func (m *DeleteDataModel) Init() tea.Cmd {
	return nil
}

func (m *DeleteDataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			if m.model.currentData != nil {
				if err := m.model.client.DeleteData(m.model.currentData.Id); err != nil {
					m.model.err = err
				} else {
					m.model.message = "Данные успешно удалены"
				}
			}
			m.model.state = StateListData
			listModel := NewListDataModel(m.model)
			return listModel, listModel.Init()
		case "n", "N", "esc", "q":
			m.model.state = StateViewData
			return NewViewDataModel(m.model), nil
		}
	}
	return m, nil
}

func (m *DeleteDataModel) View() string {
	var view []string
	view = append(view, titleStyle.Render("Удаление данных"))
	view = append(view, "")
	
	if m.model.currentData != nil {
		view = append(view, fmt.Sprintf("Вы уверены, что хотите удалить '%s'?", m.model.currentData.Name))
	} else {
		view = append(view, "Вы уверены, что хотите удалить эти данные?")
	}
	
	view = append(view, "")
	view = append(view, "y - да, n - нет")

	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, view...))
}
