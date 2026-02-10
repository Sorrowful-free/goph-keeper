package tui

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gophkeeper/gophkeeper/proto"
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
	view = append(view, fmt.Sprintf("Тип: %s", getDataTypeNameProto(data.Type)))
	view = append(view, fmt.Sprintf("ID: %s", data.Id))
	view = append(view, fmt.Sprintf("Версия: %d", data.Version))
	view = append(view, "")

	if len(data.Metadata) > 0 {
		view = append(view, "Метаданные:")
		for _, md := range data.Metadata {
			view = append(view, fmt.Sprintf("  %s: %s", md.Key, md.Value))
		}
		view = append(view, "")
	}

	// Декодируем и показываем содержимое по типу
	if content := formatDataContent(data); content != "" {
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

// formatDataContent декодирует EncryptedData и возвращает читаемый текст по типу записи
func formatDataContent(data *proto.Data) string {
	if data == nil || len(data.EncryptedData) == 0 {
		return ""
	}
	payload := data.EncryptedData

	switch data.Type {
	case proto.DataType_LOGIN_PASSWORD:
		var v struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := json.Unmarshal(payload, &v); err != nil {
			return fmt.Sprintf("  (ошибка декодирования: %v)", err)
		}
		return fmt.Sprintf("  Логин:    %s\n  Пароль:   %s", v.Login, v.Password)
	case proto.DataType_TEXT:
		var v struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(payload, &v); err != nil {
			return fmt.Sprintf("  %s", string(payload))
		}
		return "  " + v.Text
	case proto.DataType_BANK_CARD:
		var v struct {
			Number string `json:"number"`
			Expiry string `json:"expiry"`
			CVV    string `json:"cvv"`
			Holder string `json:"holder"`
		}
		if err := json.Unmarshal(payload, &v); err != nil {
			return fmt.Sprintf("  (ошибка декодирования: %v)", err)
		}
		return fmt.Sprintf("  Номер:    %s\n  Срок:     %s\n  CVV:      %s\n  Держатель: %s",
			v.Number, v.Expiry, v.CVV, v.Holder)
	case proto.DataType_BINARY:
		return "  " + string(payload)
	default:
		return "  " + string(payload)
	}
}

func getDataTypeNameProto(dt proto.DataType) string {
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
