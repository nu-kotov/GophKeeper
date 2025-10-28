package locminiio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

// BinaryStorage - структура хранилища под бинарные данные.
type BinaryStorage struct {
	stor *minio.Client
}

// NewBinaryStorage — конструктор хранилища под бинарные данные.
func NewBinaryStorage(stor *minio.Client) *BinaryStorage {
	return &BinaryStorage{stor: stor}
}

// InsertBinaryData - вставка бинарных данных в miniio.
func (bin *BinaryStorage) InsertBinaryData(ctx context.Context, bucketName string, filename string, body io.ReadCloser) (*minio.UploadInfo, error) {
	info, err := bin.stor.PutObject(ctx, bucketName, filename, body, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// SelectBinaryData - получение бинарных данных из miniio.
func (bin *BinaryStorage) SelectBinaryData(ctx context.Context, bucketName string, filename string) (*minio.Object, error) {

	file, err := bin.stor.GetObject(context.Background(), bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return file, nil
}

// DeleteBinaryData - удаление бинарных данных из miniio.
func (bin *BinaryStorage) DeleteBinaryData(ctx context.Context, bucketName string, filename string) error {
	err := bin.stor.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
