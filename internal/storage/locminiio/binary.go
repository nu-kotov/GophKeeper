package locminiio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type BinaryStorage struct {
	Stor *minio.Client
}

func (bin *BinaryStorage) InsertBinaryData(ctx context.Context, bucketName string, filename string, body io.ReadCloser) (*minio.UploadInfo, error) {
	info, err := bin.Stor.PutObject(ctx, bucketName, filename, body, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (bin *BinaryStorage) SelectBinaryData(ctx context.Context, bucketName string, filename string) (*minio.Object, error) {

	file, err := bin.Stor.GetObject(context.Background(), bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (bin *BinaryStorage) DeleteBinaryData(ctx context.Context, bucketName string, filename string) error {
	err := bin.Stor.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
