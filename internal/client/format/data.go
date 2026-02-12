package format

import (
	"time"

	"github.com/google/uuid"
	"github.com/gophkeeper/gophkeeper/proto"
)

// BuildDataForSave создаёт *proto.Data для сохранения (id, timestamps, version).
func BuildDataForSave(name string, dataType proto.DataType, encryptedPayload []byte) *proto.Data {
	now := time.Now().Unix()
	return &proto.Data{
		Id:            uuid.New().String(),
		Type:          dataType,
		Name:          name,
		EncryptedData: encryptedPayload,
		CreatedAt:     now,
		UpdatedAt:     now,
		Version:       1,
	}
}
