package locminiio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/logger"
)

func NewMiniIOConnect(conn config.MiniIOConnection) (*minio.Client, error) {

	minioClient, err := minio.New(conn.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conn.AccessKeyID, conn.SecretAccessKey, ""),
		Secure: conn.UseSSL,
	})

	if err != nil {
		logger.Log.Info(fmt.Sprintf("miniio client creation error: %s", err.Error()))
		return nil, err
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, conn.BucketName)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("bucket check error: %s", err.Error()))
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, conn.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			logger.Log.Info(fmt.Sprintf("bucket creation error: %s", err.Error()))
			return nil, err
		}

	}

	return minioClient, nil
}
