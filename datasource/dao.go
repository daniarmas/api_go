package datasource

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daniarmas/api_go/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DAO interface {
	NewObjectStorageDatasource() ObjectStorageDatasource
}

type dao struct{}

var DB *gorm.DB
var Config *utils.Config
var MinioClient *minio.Client

func NewDAO(db *gorm.DB, config *utils.Config, minio *minio.Client) DAO {
	DB = db
	Config = config
	MinioClient = minio
	return &dao{}
}

func NewMinioClient(config *utils.Config) (*minio.Client, error) {
	minioClient, err := minio.New(config.ObjectStorageServerEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.ObjectStorageServerAccessKeyId, config.ObjectStorageServerSecretAccessKey, ""),
		Secure: config.ObjectStorageServerUseSsl,
	})
	if err != nil {
		log.Fatalln(err)
	}
	businessAvatarRes, businessAvatarErr := minioClient.BucketExists(context.Background(), config.BusinessAvatarBulkName)
	if businessAvatarErr != nil {
		return nil, businessAvatarErr
	}
	if !businessAvatarRes {
		err = minioClient.MakeBucket(context.Background(), config.BusinessAvatarBulkName, minio.MakeBucketOptions{ObjectLocking: false})
		if err != nil {
			return nil, err
		}
		setPolicyErr := minioClient.SetBucketPolicy(context.Background(), config.BusinessAvatarBulkName, `{
    	"Version": "2012-10-17",
    	"Statement": [
        {
            "Sid": "PublicRead",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject",
                "s3:GetObjectVersion"
            ],
            "Resource": [
                "arn:aws:s3:::business-avatar/*"
            ]
        }
    	]
		}`)
		if setPolicyErr != nil {
			return nil, setPolicyErr
		}
	}
	itemsRes, itemsErr := minioClient.BucketExists(context.Background(), config.ItemsBulkName)
	if itemsErr != nil {
		return nil, itemsErr
	}
	if !itemsRes {
		err = minioClient.MakeBucket(context.Background(), config.ItemsBulkName, minio.MakeBucketOptions{ObjectLocking: false})
		if err != nil {
			return nil, err
		}
		setPolicyErr := minioClient.SetBucketPolicy(context.Background(), config.ItemsBulkName, `{
    	"Version": "2012-10-17",
    	"Statement": [
        {
            "Sid": "PublicRead",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject",
                "s3:GetObjectVersion"
            ],
            "Resource": [
                "arn:aws:s3:::items/*"
            ]
        }
    	]
		}`)
		if setPolicyErr != nil {
			return nil, setPolicyErr
		}
	}
	userAvatarRes, userAvatarErr := minioClient.BucketExists(context.Background(), config.UsersBulkName)
	if userAvatarErr != nil {
		return nil, userAvatarErr
	}
	if !userAvatarRes {
		err = minioClient.MakeBucket(context.Background(), config.UsersBulkName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
		if err != nil {
			return nil, err
		}
		setPolicyErr := minioClient.SetBucketPolicy(context.Background(), config.UsersBulkName, `{
    	"Version": "2012-10-17",
    	"Statement": [
        {
            "Sid": "PublicRead",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject",
                "s3:GetObjectVersion"
            ],
            "Resource": [
                "arn:aws:s3:::user-avatar/*"
            ]
        }
    	]
		}`)
		if setPolicyErr != nil {
			return nil, setPolicyErr
		}
	}
	return minioClient, nil
}

func NewConfig() (*utils.Config, error) {
	Config, err := utils.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	return &Config, nil
}

func NewDB(config *utils.Config) (*gorm.DB, error) {
	host := config.DBHost
	port := config.DBPort
	user := config.DBUser
	dbName := config.DBDatabase
	password := config.DBPassword

	// Starting a database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func (d *dao) NewObjectStorageDatasource() ObjectStorageDatasource {
	return &objectStorageDatasource{Minio: MinioClient}
}