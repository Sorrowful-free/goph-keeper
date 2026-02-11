package tui

import (
	"iter"
	"slices"

	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/gophkeeper/gophkeeper/internal/client/format"
	"github.com/gophkeeper/gophkeeper/proto"
)

const (
	addDataStepNameType = 0
	addDataStepFields   = 1
)

// AddDataModel представляет модель добавления данных
type AddDataModel struct {
	model      *Model
	step       int
	nameInput  textinput.Model
	typeSelect int
	types      []string
	focused    int
	err        error

	// Поля шага 2 (по типу)
	loginInput   textinput.Model
	passwordInput textinput.Model
	textInput    textinput.Model
	cardNumber   textinput.Model
	cardExpiry   textinput.Model
	cardCVV      textinput.Model
	cardHolder   textinput.Model
	binaryInput  textinput.Model
	fieldFocus   int
}

func newLoginInput() textinput.Model {
	t := textinput.New()
	t.Placeholder = "Логин"
	t.CharLimit = 200
	t.Width = 38
	return t
}

func newPasswordInput() textinput.Model {
	t := textinput.New()
	t.Placeholder = "Пароль"
	t.CharLimit = 200
	t.Width = 38
	t.EchoMode = textinput.EchoPassword
	t.EchoCharacter = '•'
	return t
}

func newTextInput(placeholder string) textinput.Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.CharLimit = 2000
	t.Width = 38
	return t
}

func NewAddDataModel(m *Model) *AddDataModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Название"
	nameInput.Focus()
	nameInput.CharLimit = 100
	nameInput.Width = 38

	types := format.DataTypeDisplayNames()

	return &AddDataModel{
		model:       m,
		step:        addDataStepNameType,
		nameInput:   nameInput,
		typeSelect:  0,
		types:       types,
		focused:     0,
		loginInput:  newLoginInput(),
		passwordInput: newPasswordInput(),
		textInput:   newTextInput("Содержимое"),
		cardNumber:  newTextInput("Номер карты"),
		cardExpiry:  newTextInput("Срок (MM/YY)"),
		cardCVV:     newTextInput("CVV"),
		cardHolder:  newTextInput("Держатель карты"),
		binaryInput: newTextInput("Данные (base64 или текст)"),
		fieldFocus:  0,
	}
}

func (m *AddDataModel) dataType() proto.DataType {
	return format.DataTypeFromIndex(m.typeSelect)
}

func (m *AddDataModel) fieldCount() int {
	return format.FieldCount(m.dataType())
}

func (m *AddDataModel) focusFirstField() {
	m.fieldFocus = 0
	m.loginInput.Blur()
	m.passwordInput.Blur()
	m.textInput.Blur()
	m.cardNumber.Blur()
	m.cardExpiry.Blur()
	m.cardCVV.Blur()
	m.cardHolder.Blur()
	m.binaryInput.Blur()
	switch m.dataType() {
	case proto.DataType_LOGIN_PASSWORD:
		m.loginInput.Focus()
	case proto.DataType_TEXT:
		m.textInput.Focus()
	case proto.DataType_BINARY:
		m.binaryInput.Focus()
	case proto.DataType_BANK_CARD:
		m.cardNumber.Focus()
	}
}

func (m *AddDataModel) moveFocusFields(delta int) tea.Cmd {
	n := m.fieldCount()
	if n == 0 {
		return nil
	}
	m.fieldFocus += delta
	if m.fieldFocus < 0 {
		m.fieldFocus = 0
	}
	if m.fieldFocus >= n {
		m.fieldFocus = n - 1
	}
	m.loginInput.Blur()
	m.passwordInput.Blur()
	m.textInput.Blur()
	m.cardNumber.Blur()
	m.cardExpiry.Blur()
	m.cardCVV.Blur()
	m.cardHolder.Blur()
	m.binaryInput.Blur()
	var focused textinput.Model
	switch m.dataType() {
	case proto.DataType_LOGIN_PASSWORD:
		if m.fieldFocus == 0 {
			focused = m.loginInput
			m.loginInput.Focus()
		} else {
			focused = m.passwordInput
			m.passwordInput.Focus()
		}
	case proto.DataType_TEXT:
		focused = m.textInput
		m.textInput.Focus()
	case proto.DataType_BINARY:
		focused = m.binaryInput
		m.binaryInput.Focus()
	case proto.DataType_BANK_CARD:
		switch m.fieldFocus {
		case 0:
			focused = m.cardNumber
			m.cardNumber.Focus()
		case 1:
			focused = m.cardExpiry
			m.cardExpiry.Focus()
		case 2:
			focused = m.cardCVV
			m.cardCVV.Focus()
		case 3:
			focused = m.cardHolder
			m.cardHolder.Focus()
		}
	}
	_ = focused
	return nil
}

func (m *AddDataModel) getFocusedInput() *textinput.Model {
	switch m.dataType() {
	case proto.DataType_LOGIN_PASSWORD:
		if m.fieldFocus == 0 {
			return &m.loginInput
		}
		return &m.passwordInput
	case proto.DataType_TEXT:
		return &m.textInput
	case proto.DataType_BINARY:
		return &m.binaryInput
	case proto.DataType_BANK_CARD:
		switch m.fieldFocus {
		case 0:
			return &m.cardNumber
		case 1:
			return &m.cardExpiry
		case 2:
			return &m.cardCVV
		case 3:
			return &m.cardHolder
		}
	}
	return &m.textInput
}

func (m *AddDataModel) buildEncryptedData() ([]byte, error) {
	fields := m.collectFieldValues()
	return format.BuildPayload(m.dataType(), fields)
}

func (m *AddDataModel) collectFieldValues() map[string]string {
	fields := make(map[string]string)
	dt := m.dataType()
	switch dt {
	case proto.DataType_LOGIN_PASSWORD:
		fields[format.FieldLogin] = m.loginInput.Value()
		fields[format.FieldPassword] = m.passwordInput.Value()
	case proto.DataType_TEXT:
		fields[format.FieldText] = m.textInput.Value()
	case proto.DataType_BINARY:
		fields[format.FieldBinary] = m.binaryInput.Value()
	case proto.DataType_BANK_CARD:
		fields[format.FieldNumber] = m.cardNumber.Value()
		fields[format.FieldExpiry] = m.cardExpiry.Value()
		fields[format.FieldCVV] = m.cardCVV.Value()
		fields[format.FieldHolder] = m.cardHolder.Value()
	}
	return fields
}

func (m *AddDataModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AddDataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if m.step == addDataStepFields {
				m.step = addDataStepNameType
				m.focusFirstField()
				m.err = nil
				return m, nil
			}
			m.model.state = StateMainMenu
			return NewMainMenuModel(m.model), nil
		case "tab", "down":
			if m.step == addDataStepNameType {
				if m.focused == 0 {
					m.focused = 1
					m.nameInput.Blur()
				} else {
					m.typeSelect = (m.typeSelect + 1) % len(m.types)
				}
			} else {
				m.moveFocusFields(1)
			}
			return m, nil
		case "up":
			if m.step == addDataStepNameType {
				if m.focused == 1 {
					if m.typeSelect > 0 {
						m.typeSelect--
					} else {
						m.focused = 0
						m.nameInput.Focus()
					}
				}
			} else {
				m.moveFocusFields(-1)
			}
			return m, nil
		case "enter":
			if m.step == addDataStepNameType {
				name := m.nameInput.Value()
				if name == "" {
					m.err = fmt.Errorf("название обязательно")
					return m, nil
				}
				m.err = nil
				m.step = addDataStepFields
				m.focusFirstField()
			} else {
				return m.handleSubmit()
			}
			return m, nil
		}
	case error:
		m.err = msg
		return m, nil
	}

	if m.step == addDataStepNameType && m.focused == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
		return m, cmd
	}
	if m.step == addDataStepFields {
		inp := m.getFocusedInput()
		*inp, cmd = inp.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *AddDataModel) handleSubmit() (tea.Model, tea.Cmd) {
	name := m.nameInput.Value()
	if name == "" {
		m.err = fmt.Errorf("название обязательно")
		return m, nil
	}

	payload, err := m.buildEncryptedData()
	if err != nil {
		m.err = err
		return m, nil
	}

	data := format.BuildDataForSave(name, m.dataType(), payload)
	_, _, err = m.model.client.SaveData(data)
	if err != nil {
		m.err = err
		return m, nil
	}

	m.model.message = "Данные успешно сохранены"
	m.model.state = StateMainMenu
	return NewMainMenuModel(m.model), nil
}

func (m *AddDataModel) View() string {
	var view []string
	view = append(view, titleStyle.Render("Добавление данных"))
	view = append(view, "")

	if m.step == addDataStepNameType {
		if m.focused == 0 {
			view = append(view, focusedStyle.Render(m.nameInput.View()))
		} else {
			view = append(view, inputStyle.Render(m.nameInput.View()))
		}
		view = append(view, "")
		view = append(view, "Тип данных:")
		view = append(view, slices.Collect(addDataTypesToLinesSeq(m.types, m.typeSelect, m.focused == 1))...)
		view = append(view, "")
		view = append(view, "Tab/↓ — следующий, Enter — далее или сохранить, Esc — назад/отмена")
	} else {
		view = append(view, inputStyle.Render("Название: "+m.nameInput.Value()))
		view = append(view, "")
		view = append(view, "Тип: "+m.types[m.typeSelect])
		view = append(view, "")

		switch m.dataType() {
		case proto.DataType_LOGIN_PASSWORD:
			if m.fieldFocus == 0 {
				view = append(view, focusedStyle.Render(m.loginInput.View()))
			} else {
				view = append(view, inputStyle.Render(m.loginInput.View()))
			}
			view = append(view, "")
			if m.fieldFocus == 1 {
				view = append(view, focusedStyle.Render(m.passwordInput.View()))
			} else {
				view = append(view, inputStyle.Render(m.passwordInput.View()))
			}
		case proto.DataType_TEXT:
			if m.fieldFocus == 0 {
				view = append(view, focusedStyle.Render(m.textInput.View()))
			} else {
				view = append(view, inputStyle.Render(m.textInput.View()))
			}
		case proto.DataType_BINARY:
			if m.fieldFocus == 0 {
				view = append(view, focusedStyle.Render(m.binaryInput.View()))
			} else {
				view = append(view, inputStyle.Render(m.binaryInput.View()))
			}
		case proto.DataType_BANK_CARD:
			style := inputStyle
			if m.fieldFocus == 0 {
				style = focusedStyle
			}
			view = append(view, style.Render(m.cardNumber.View()))
			view = append(view, "")
			style = inputStyle
			if m.fieldFocus == 1 {
				style = focusedStyle
			}
			view = append(view, style.Render(m.cardExpiry.View()))
			view = append(view, "")
			style = inputStyle
			if m.fieldFocus == 2 {
				style = focusedStyle
			}
			view = append(view, style.Render(m.cardCVV.View()))
			view = append(view, "")
			style = inputStyle
			if m.fieldFocus == 3 {
				style = focusedStyle
			}
			view = append(view, style.Render(m.cardHolder.View()))
		}
		view = append(view, "")
		view = append(view, "Tab/↓ — следующее поле, Enter — сохранить, Esc — назад")
	}

	if m.err != nil {
		view = append(view, "")
		view = append(view, errorStyle.Render(fmt.Sprintf("Ошибка: %v", m.err)))
	}

	view = append(view, "")

	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, view...))
}

// addDataTypesToLinesSeq возвращает итератор строк списка типов (выбранный/обычный стиль)
func addDataTypesToLinesSeq(types []string, typeSelect int, typeFocused bool) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i, t := range types {
			var line string
			if typeFocused && i == typeSelect {
				line = selectedMenuItemStyle.Render("▶ " + t)
			} else {
				line = menuItemStyle.Render("  " + t)
			}
			if !yield(line) {
				return
			}
		}
	}
}
