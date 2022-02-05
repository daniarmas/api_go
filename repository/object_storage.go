package repository

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageRepository interface {
	BucketExists(ctx context.Context, bucketName string) (*bool, error)
	ObjectExists(ctx context.Context, bucketName string, objectName string) (*bool, error)
	CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (*minio.UploadInfo, error)
}

type objectStorageRepository struct {
}

func (m *objectStorageRepository) BucketExists(ctx context.Context, bucketName string) (*bool, error) {
	result, err := Datasource.NewObjectStorageDatasource().BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *objectStorageRepository) ObjectExists(ctx context.Context, bucketName string, objectName string) (*bool, error) {
	result, err := Datasource.NewObjectStorageDatasource().ObjectExists(ctx, bucketName, objectName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *objectStorageRepository) CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (*minio.UploadInfo, error) {
	result, err := Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), dst, src)
	if err != nil {
		return nil, err
	}
	return result, nil
}