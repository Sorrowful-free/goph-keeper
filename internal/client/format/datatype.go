package format

import "github.com/gophkeeper/gophkeeper/proto"

// DataTypeDisplayName возвращает человекочитаемое имя типа данных.
func DataTypeDisplayName(dt proto.DataType) string {
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

// DataTypeDisplayNames возвращает список отображаемых имён типов в порядке,
// соответствующем индексу выбора в UI (Логин/Пароль, Текст, Бинарные данные, Банковская карта).
func DataTypeDisplayNames() []string {
	return []string{
		"Логин/Пароль",
		"Текст",
		"Бинарные данные",
		"Банковская карта",
	}
}

// DataTypeFromIndex возвращает proto.DataType по индексу в списке типов (0..3).
func DataTypeFromIndex(i int) proto.DataType {
	switch i {
	case 0:
		return proto.DataType_LOGIN_PASSWORD
	case 1:
		return proto.DataType_TEXT
	case 2:
		return proto.DataType_BINARY
	case 3:
		return proto.DataType_BANK_CARD
	}
	return proto.DataType_UNKNOWN
}

// FieldCount возвращает количество полей ввода для типа данных (для формы добавления).
func FieldCount(dt proto.DataType) int {
	switch dt {
	case proto.DataType_LOGIN_PASSWORD:
		return 2
	case proto.DataType_TEXT, proto.DataType_BINARY:
		return 1
	case proto.DataType_BANK_CARD:
		return 4
	}
	return 0
}
