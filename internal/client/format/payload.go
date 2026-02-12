package format

import (
	"encoding/json"

	"github.com/gophkeeper/gophkeeper/proto"
)

// Ключи полей для BuildPayload (форма добавления).
const (
	FieldLogin    = "login"
	FieldPassword = "password"
	FieldText     = "text"
	FieldBinary   = "binary"
	FieldNumber   = "number"
	FieldExpiry   = "expiry"
	FieldCVV      = "cvv"
	FieldHolder   = "holder"
)

// BuildPayload собирает EncryptedData из полей формы по типу данных.
// fields — значения полей по ключам (FieldLogin, FieldPassword, FieldText и т.д.).
func BuildPayload(dataType proto.DataType, fields map[string]string) ([]byte, error) {
	switch dataType {
	case proto.DataType_LOGIN_PASSWORD:
		return json.Marshal(map[string]string{
			"login":    fields[FieldLogin],
			"password": fields[FieldPassword],
		})
	case proto.DataType_TEXT:
		return json.Marshal(map[string]string{"text": fields[FieldText]})
	case proto.DataType_BINARY:
		return []byte(fields[FieldBinary]), nil
	case proto.DataType_BANK_CARD:
		return json.Marshal(map[string]string{
			"number": fields[FieldNumber],
			"expiry": fields[FieldExpiry],
			"cvv":    fields[FieldCVV],
			"holder": fields[FieldHolder],
		})
	}
	return []byte{}, nil
}
