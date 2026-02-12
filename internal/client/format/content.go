package format

import (
	"encoding/json"
	"fmt"

	"github.com/gophkeeper/gophkeeper/proto"
)

// DataContentToDisplayString декодирует EncryptedData и возвращает читаемый текст по типу записи.
func DataContentToDisplayString(data *proto.Data) string {
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
