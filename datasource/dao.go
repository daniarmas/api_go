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
	NewJwtTokenDatasource() JwtTokenDatasource
	NewVerificationCodeDatasource() VerificationCodeDatasource
	NewUserDatasource() UserDatasource
	NewRefreshTokenDatasource() RefreshTokenDatasource
	NewAuthorizationTokenDatasource() AuthorizationTokenDatasource
	NewItemDatasource() ItemDatasource
	NewItemCategoryDatasource() ItemCategoryDatasource
	NewDeviceDatasource() DeviceDatasource
	NewBusinessDatasource() BusinessDatasource
	NewBannedUserDatasource() BannedUserDatasource
	NewBusinessUserDatasource() BusinessUserDatasource
	NewBannedDeviceDatasource() BannedDeviceDatasource
	NewCartItemDatasource() CartItemDatasource
	NewPermissionDatasource() PermissionDatasource
	NewProvinceDatasource() ProvinceDatasource
	NewMunicipalityDatasource() MunicipalityDatasource
	NewUnionBusinessAndMunicipalityDatasource() UnionBusinessAndMunicipalityDatasource
	NewOrderDatasource() OrderDatasource
	NewOrderedItemDatasource() OrderedItemDatasource
	NewUnionOrderAndOrderedItemDatasource() UnionOrderAndOrderedItemDatasource
	NewBannedAppDatasource() BannedAppDatasource
	NewBusinessScheduleDatasource() BusinessScheduleDatasource
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
	var secure bool
	if config.ObjectStorageServerUseSsl == "true" {
		secure = true
	}
	minioClient, err := minio.New(config.ObjectStorageServerEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.ObjectStorageServerAccessKeyId, config.ObjectStorageServerSecretAccessKey, ""),
		Secure: secure,
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
	itemsDeletedRes, itemsDeletedErr := minioClient.BucketExists(context.Background(), config.ItemsDeletedBulkName)
	if itemsDeletedErr != nil {
		return nil, itemsDeletedErr
	}
	if !itemsDeletedRes {
		err = minioClient.MakeBucket(context.Background(), config.ItemsDeletedBulkName, minio.MakeBucketOptions{ObjectLocking: false})
		if err != nil {
			return nil, err
		}
	}
	userAvatarRes, userAvatarErr := minioClient.BucketExists(context.Background(), config.UsersBulkName)
	if userAvatarErr != nil {
		return nil, userAvatarErr
	}
	if !userAvatarRes {
		err = minioClient.MakeBucket(context.Background(), config.UsersBulkName, minio.MakeBucketOptions{ObjectLocking: false})
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
	usersDeletedRes, usersDeletedErr := minioClient.BucketExists(context.Background(), config.UsersDeletedBulkName)
	if usersDeletedErr != nil {
		return nil, usersDeletedErr
	}
	if !usersDeletedRes {
		err = minioClient.MakeBucket(context.Background(), config.UsersDeletedBulkName, minio.MakeBucketOptions{ObjectLocking: false})
		if err != nil {
			return nil, err
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

func (d *dao) NewOrderedItemDatasource() OrderedItemDatasource {
	return &orderedItemDatasource{}
}

func (d *dao) NewBannedAppDatasource() BannedAppDatasource {
	return &bannedAppDatasource{}
}

func (d *dao) NewBusinessScheduleDatasource() BusinessScheduleDatasource {
	return &businessScheduleDatasource{}
}

func (d *dao) NewUnionOrderAndOrderedItemDatasource() UnionOrderAndOrderedItemDatasource {
	return &unionOrderAndOrderedItemDatasource{}
}

func (d *dao) NewJwtTokenDatasource() JwtTokenDatasource {
	return &jwtTokenDatasource{}
}

func (d *dao) NewBusinessUserDatasource() BusinessUserDatasource {
	return &businessUserDatasource{}
}

func (d *dao) NewVerificationCodeDatasource() VerificationCodeDatasource {
	return &verificationCodeDatasource{}
}

func (d *dao) NewUserDatasource() UserDatasource {
	return &userDatasource{}
}

func (d *dao) NewRefreshTokenDatasource() RefreshTokenDatasource {
	return &refreshTokenDatasource{}
}

func (d *dao) NewUnionBusinessAndMunicipalityDatasource() UnionBusinessAndMunicipalityDatasource {
	return &unionBusinessAndMunicipalityDatasource{}
}

func (d *dao) NewAuthorizationTokenDatasource() AuthorizationTokenDatasource {
	return &authorizationTokenDatasource{}
}

func (d *dao) NewItemDatasource() ItemDatasource {
	return &itemDatasource{}
}

func (d *dao) NewProvinceDatasource() ProvinceDatasource {
	return &provinceDatasource{}
}

func (d *dao) NewMunicipalityDatasource() MunicipalityDatasource {
	return &municipalityDatasource{}
}

func (d *dao) NewItemCategoryDatasource() ItemCategoryDatasource {
	return &itemCategoryDatasource{}
}

func (d *dao) NewPermissionDatasource() PermissionDatasource {
	return &permissionDatasource{}
}

func (d *dao) NewOrderDatasource() OrderDatasource {
	return &orderDatasource{}
}

func (d *dao) NewDeviceDatasource() DeviceDatasource {
	return &deviceDatasource{}
}

func (d *dao) NewBusinessDatasource() BusinessDatasource {
	return &businessDatasource{}
}

func (d *dao) NewBannedUserDatasource() BannedUserDatasource {
	return &bannedUserDatasource{}
}

func (d *dao) NewBannedDeviceDatasource() BannedDeviceDatasource {
	return &bannedDeviceDatasource{}
}

func (d *dao) NewCartItemDatasource() CartItemDatasource {
	return &cartItemDatasource{}
}
