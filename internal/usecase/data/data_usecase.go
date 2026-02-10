package data

import (
	"github.com/gophkeeper/gophkeeper/internal/domain/repository"
)

// DataUseCase объединяет сценарии работы с данными пользователя
type DataUseCase struct {
	dataRepo repository.DataRepository
}

// NewDataUseCase создаёт use case данных
func NewDataUseCase(dataRepo repository.DataRepository) *DataUseCase {
	return &DataUseCase{
		dataRepo: dataRepo,
	}
}
