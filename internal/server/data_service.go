package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gophkeeper/gophkeeper/internal/crypto"
	"github.com/gophkeeper/gophkeeper/internal/models"
	"github.com/gophkeeper/gophkeeper/internal/usecase/data"
	"github.com/gophkeeper/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// DataService реализует gRPC-сервис работы с данными (delivery layer)
type DataService struct {
	proto.UnimplementedDataServiceServer
	dataUC *data.DataUseCase
}

// NewDataService создаёт новый сервис данных
func NewDataService(dataUC *data.DataUseCase) *DataService {
	return &DataService{
		dataUC: dataUC,
	}
}

// getUserIDFromContext извлекает user ID из контекста (из JWT токена)
func (s *DataService) getUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no metadata")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", status.Error(codes.Unauthenticated, "no authorization token")
	}

	token := tokens[0]
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := crypto.ValidateToken(token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims.UserID, nil
}

// convertProtoDataType конвертирует proto DataType в models.DataType
func convertProtoDataType(protoType proto.DataType) models.DataType {
	switch protoType {
	case proto.DataType_LOGIN_PASSWORD:
		return models.DataTypeLoginPassword
	case proto.DataType_TEXT:
		return models.DataTypeText
	case proto.DataType_BINARY:
		return models.DataTypeBinary
	case proto.DataType_BANK_CARD:
		return models.DataTypeBankCard
	default:
		return models.DataTypeText
	}
}

// convertModelsDataType конвертирует models.DataType в proto.DataType
func convertModelsDataType(modelType models.DataType) proto.DataType {
	switch modelType {
	case models.DataTypeLoginPassword:
		return proto.DataType_LOGIN_PASSWORD
	case models.DataTypeText:
		return proto.DataType_TEXT
	case models.DataTypeBinary:
		return proto.DataType_BINARY
	case models.DataTypeBankCard:
		return proto.DataType_BANK_CARD
	default:
		return proto.DataType_UNKNOWN
	}
}

// convertProtoDataToModel конвертирует proto.Data в models.Data
func convertProtoDataToModel(protoData *proto.Data, userID string) (*models.Data, error) {
	metadataJSON := "{}"
	if len(protoData.Metadata) > 0 {
		metadataBytes, err := json.Marshal(protoData.Metadata)
		if err != nil {
			return nil, err
		}
		metadataJSON = string(metadataBytes)
	}

	createdAt := time.Now()
	updatedAt := time.Now()
	if protoData.CreatedAt > 0 {
		createdAt = time.Unix(protoData.CreatedAt, 0)
	}
	if protoData.UpdatedAt > 0 {
		updatedAt = time.Unix(protoData.UpdatedAt, 0)
	}

	encryptedData := protoData.EncryptedData
	if encryptedData == nil {
		encryptedData = []byte{}
	}

	return &models.Data{
		ID:            protoData.Id,
		UserID:        userID,
		Type:          convertProtoDataType(protoData.Type),
		Name:          protoData.Name,
		EncryptedData: encryptedData,
		Metadata:      metadataJSON,
		Version:       protoData.Version,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

// convertModelDataToProto конвертирует models.Data в proto.Data
func convertModelDataToProto(modelData *models.Data) (*proto.Data, error) {
	var metadataItems []*proto.Metadata
	if modelData.Metadata != "" && modelData.Metadata != "{}" {
		var metadataList []models.MetadataItem
		if err := json.Unmarshal([]byte(modelData.Metadata), &metadataList); err == nil {
			for _, item := range metadataList {
				metadataItems = append(metadataItems, &proto.Metadata{
					Key:   item.Key,
					Value: item.Value,
				})
			}
		}
	}

	return &proto.Data{
		Id:            modelData.ID,
		Type:          convertModelsDataType(modelData.Type),
		Name:          modelData.Name,
		EncryptedData: modelData.EncryptedData,
		Metadata:      metadataItems,
		CreatedAt:     modelData.CreatedAt.Unix(),
		UpdatedAt:     modelData.UpdatedAt.Unix(),
		Version:       modelData.Version,
	}, nil
}

// SaveData сохраняет данные пользователя
func (s *DataService) SaveData(ctx context.Context, req *proto.SaveDataRequest) (*proto.SaveDataResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return &proto.SaveDataResponse{
			Success: false,
			Message: "authentication required",
		}, err
	}

	if req.Data == nil {
		return &proto.SaveDataResponse{
			Success: false,
			Message: "data is required",
		}, status.Error(codes.InvalidArgument, "data is required")
	}

	modelData, err := convertProtoDataToModel(req.Data, userID)
	if err != nil {
		return &proto.SaveDataResponse{
			Success: false,
			Message: fmt.Sprintf("error converting data: %v", err),
		}, status.Error(codes.InvalidArgument, "invalid data format")
	}

	out, err := s.dataUC.SaveData(ctx, data.SaveDataInput{
		UserID: userID,
		Data:   modelData,
	})
	if err != nil {
		return &proto.SaveDataResponse{
			Success: false,
			Message: fmt.Sprintf("error saving data: %v", err),
		}, status.Error(codes.Internal, "internal error")
	}

	return &proto.SaveDataResponse{
		Success: true,
		Message: "data saved successfully",
		DataId:  out.DataID,
		Version: out.Version,
	}, nil
}

// GetData получает данные по ID
func (s *DataService) GetData(ctx context.Context, req *proto.GetDataRequest) (*proto.GetDataResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return &proto.GetDataResponse{
			Success: false,
			Message: "authentication required",
		}, err
	}

	if req.DataId == "" {
		return &proto.GetDataResponse{
			Success: false,
			Message: "data_id is required",
		}, status.Error(codes.InvalidArgument, "data_id is required")
	}

	out, err := s.dataUC.GetData(ctx, data.GetDataInput{
		UserID: userID,
		DataID: req.DataId,
	})
	if err != nil {
		if errors.Is(err, data.ErrDataNotFound) {
			return &proto.GetDataResponse{
				Success: false,
				Message: "data not found",
			}, status.Error(codes.NotFound, "data not found")
		}
		return &proto.GetDataResponse{
			Success: false,
			Message: "data not found",
		}, status.Error(codes.NotFound, "data not found")
	}

	protoData, err := convertModelDataToProto(out.Data)
	if err != nil {
		return &proto.GetDataResponse{
			Success: false,
			Message: "error converting data",
		}, status.Error(codes.Internal, "internal error")
	}

	return &proto.GetDataResponse{
		Success: true,
		Message: "data retrieved successfully",
		Data:    protoData,
	}, nil
}

// ListData получает список данных пользователя
func (s *DataService) ListData(ctx context.Context, req *proto.ListDataRequest) (*proto.ListDataResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return &proto.ListDataResponse{
			Success: false,
			Message: "authentication required",
		}, err
	}

	var dataType models.DataType
	if req.Type != proto.DataType_UNKNOWN {
		dataType = convertProtoDataType(req.Type)
	}

	out, err := s.dataUC.ListData(ctx, data.ListDataInput{
		UserID:   userID,
		DataType: dataType,
	})
	if err != nil {
		return &proto.ListDataResponse{
			Success: false,
			Message: fmt.Sprintf("error listing data: %v", err),
		}, status.Error(codes.Internal, "internal error")
	}

	protoDataList := make([]*proto.Data, 0, len(out.Items))
	for _, d := range out.Items {
		protoData, err := convertModelDataToProto(d)
		if err != nil {
			continue
		}
		protoDataList = append(protoDataList, protoData)
	}

	return &proto.ListDataResponse{
		Success: true,
		Message: "data listed successfully",
		Data:    protoDataList,
	}, nil
}

// DeleteData удаляет данные
func (s *DataService) DeleteData(ctx context.Context, req *proto.DeleteDataRequest) (*proto.DeleteDataResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return &proto.DeleteDataResponse{
			Success: false,
			Message: "authentication required",
		}, err
	}

	if req.DataId == "" {
		return &proto.DeleteDataResponse{
			Success: false,
			Message: "data_id is required",
		}, status.Error(codes.InvalidArgument, "data_id is required")
	}

	if err := s.dataUC.DeleteData(ctx, data.DeleteDataInput{
		UserID: userID,
		DataID: req.DataId,
	}); err != nil {
		return &proto.DeleteDataResponse{
			Success: false,
			Message: "error deleting data",
		}, status.Error(codes.Internal, "internal error")
	}

	return &proto.DeleteDataResponse{
		Success: true,
		Message: "data deleted successfully",
	}, nil
}

// SyncData синхронизирует данные
func (s *DataService) SyncData(ctx context.Context, req *proto.SyncDataRequest) (*proto.SyncDataResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return &proto.SyncDataResponse{
			Success: false,
			Message: "authentication required",
		}, err
	}

	var since time.Time
	if req.LastSyncTime > 0 {
		since = time.Unix(req.LastSyncTime, 0)
	}

	out, err := s.dataUC.SyncData(ctx, data.SyncDataInput{
		UserID:       userID,
		LastSyncTime: since,
	})
	if err != nil {
		return &proto.SyncDataResponse{
			Success: false,
			Message: fmt.Sprintf("error syncing data: %v", err),
		}, status.Error(codes.Internal, "internal error")
	}

	protoDataList := make([]*proto.Data, 0, len(out.Items))
	for _, d := range out.Items {
		protoData, err := convertModelDataToProto(d)
		if err != nil {
			continue
		}
		protoDataList = append(protoDataList, protoData)
	}

	return &proto.SyncDataResponse{
		Success:  true,
		Message:  "data synced successfully",
		Data:     protoDataList,
		SyncTime: out.SyncTime.Unix(),
	}, nil
}
