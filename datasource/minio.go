package datasource

import (
	"context"
	"errors"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageDatasource interface {
	BucketExists(ctx context.Context, bucketName string) (*bool, error)
	ObjectExists(ctx context.Context, bucketName string, objectName string) (*bool, error)
	CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (*minio.UploadInfo, error)
}

type objectStorageDatasource struct {
	Minio *minio.Client
}

func (m *objectStorageDatasource) BucketExists(ctx context.Context, bucketName string) (*bool, error) {
	result, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// https://github.com/minio/minio-go/issues/1082
func (m *objectStorageDatasource) ObjectExists(ctx context.Context, bucketName string, objectName string) (*bool, error) {
	var result = true
	_, err := m.Minio.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "AccessDenied" {
			return nil, errors.New("PathInsufficientPermission")
		}
		if errResponse.Code == "NoSuchBucket" {
			return nil, errors.New("BucketDoesNotExist")
		}
		if errResponse.Code == "InvalidBucketName" {
			return nil, errors.New("BucketInvalid")
		}
		if errResponse.Code == "NoSuchKey" {
			return nil, errors.New("ObjectMissing")
		}
		return nil, err
	}
	return &result, nil
}

func (m *objectStorageDatasource) CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (*minio.UploadInfo, error) {
	result, err := m.Minio.CopyObject(context.Background(), dst, src)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
