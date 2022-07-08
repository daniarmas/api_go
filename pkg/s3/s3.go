package s3

import (
	"github.com/daniarmas/api_go/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// New - instantiate minio client with config
func New(config *config.Config) (*minio.Client, error) {
	var secure bool
	if config.ObjectStorageServerUseSsl == "true" {
		secure = true
	}
	minioClient, err := minio.New(config.ObjectStorageServerEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.ObjectStorageServerAccessKeyId, config.ObjectStorageServerSecretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}
