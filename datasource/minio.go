package datasource

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageDatasource interface {
	BucketExists(ctx context.Context, bucketName string) (bool, error)
}

type objectStorageDatasource struct {
	Minio *minio.Client
}

func (m *objectStorageDatasource) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	result, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return result, nil
}
