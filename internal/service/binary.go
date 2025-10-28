package service

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

// BinaryStorage - интерфейс хранилища для работы с бинарными данными.
type BinaryStorage interface {
	InsertBinaryData(context.Context, string, string, io.ReadCloser) (*minio.UploadInfo, error)
	SelectBinaryData(context.Context, string, string) (*minio.Object, error)
	DeleteBinaryData(context.Context, string, string) error
}

// BinaryService — структура сервиса для работы с бинарными данными.
type BinaryService struct {
	storage BinaryStorage
}

// NewBinaryService — конструктор сервиса для работы с бинарными данными.
func NewBinaryService(storage BinaryStorage) *BinaryService {
	return &BinaryService{storage: storage}
}

// AddBinary — добавление бинарных данных в miniio.
func (s *BinaryService) AddBinary(ctx context.Context, bucketName string, miniioFilename string, body io.ReadCloser) (*minio.UploadInfo, error) {
	return s.storage.InsertBinaryData(ctx, bucketName, miniioFilename, body)
}

// GetBinary — получение бинарных данных из miniio.
func (s *BinaryService) GetBinary(ctx context.Context, bucketName, miniioFilename string) (*minio.Object, error) {
	return s.storage.SelectBinaryData(ctx, bucketName, miniioFilename)
}

// DeleteBinary — удаление бинарных данных из miniio.
func (s *BinaryService) DeleteBinary(ctx context.Context, userID, dataID string) error {
	return s.storage.DeleteBinaryData(ctx, userID, dataID)
}
