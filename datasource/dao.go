package datasource

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/daniarmas/api_go/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type JsonWebTokenMetadata struct {
	TokenId *uuid.UUID
	Token   *string
}

type DAO interface {
	NewObjectStorageDatasource() ObjectStorageDatasource
	NewJwtTokenDatasource() JwtTokenDatasource
	NewVerificationCodeDatasource() VerificationCodeDatasource
	NewUserDatasource() UserDatasource
	NewRefreshTokenDatasource() RefreshTokenDatasource
	NewAuthorizationTokenDatasource() AuthorizationTokenDatasource
	NewItemDatasource() ItemDatasource
	NewBusinessCollectionDatasource() BusinessCollectionDatasource
	NewDeviceDatasource() DeviceDatasource
	NewBusinessDatasource() BusinessDatasource
	NewBannedUserDatasource() BannedUserDatasource
	NewBusinessUserDatasource() BusinessUserDatasource
	NewBannedDeviceDatasource() BannedDeviceDatasource
	NewCartItemDatasource() CartItemDatasource
	NewPermissionDatasource() PermissionDatasource
	NewUserPermissionDatasource() UserPermissionDatasource
	NewProvinceDatasource() ProvinceDatasource
	NewMunicipalityDatasource() MunicipalityDatasource
	NewUnionBusinessAndMunicipalityDatasource() UnionBusinessAndMunicipalityDatasource
	NewOrderDatasource() OrderDatasource
	NewOrderedItemDatasource() OrderedItemDatasource
	NewUnionOrderAndOrderedItemDatasource() UnionOrderAndOrderedItemDatasource
	NewApplicationDatasource() ApplicationDatasource
	NewBusinessScheduleDatasource() BusinessScheduleDatasource
	NewOrderLifecycleDatasource() OrderLifecycleDatasource
	NewBusinessCategoryDatasource() BusinessCategoryDatasource
	NewBusinessAnalyticsDatasource() BusinessAnalyticsDatasource
	NewItemAnalyticsDatasource() ItemAnalyticsDatasource
	NewUserAddressDatasource() UserAddressDatasource
	NewPartnerApplicationDatasource() PartnerApplicationDatasource
	NewBusinessRoleDatasource() BusinessRoleDatasource
	NewUnionBusinessRoleAndPermissionDatasource() UnionBusinessRoleAndPermissionDatasource
	NewUnionBusinessRoleAndUserDatasource() UnionBusinessRoleAndUserDatasource
}

type dao struct{}

var Connection *gorm.DB
var Config *config.Config
var MinioClient *minio.Client

func NewDAO(db *gorm.DB, config *config.Config, minio *minio.Client) DAO {
	Connection = db
	Config = config
	MinioClient = minio
	return &dao{}
}

func NewMinioClient(config *config.Config) (*minio.Client, error) {
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
	businessAvatarDeletedRes, businessAvatarDeletedErr := minioClient.BucketExists(context.Background(), config.BusinessAvatarDeletedBulkName)
	if businessAvatarDeletedErr != nil {
		return nil, businessAvatarDeletedErr
	}
	if !businessAvatarDeletedRes {
		err = minioClient.MakeBucket(context.Background(), config.BusinessAvatarDeletedBulkName, minio.MakeBucketOptions{ObjectLocking: false})
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

func (d *dao) NewObjectStorageDatasource() ObjectStorageDatasource {
	return &objectStorageDatasource{Minio: MinioClient}
}

func (d *dao) NewBusinessCategoryDatasource() BusinessCategoryDatasource {
	return &businessCategoryDatasource{}
}

func (d *dao) NewOrderedItemDatasource() OrderedItemDatasource {
	return &orderedItemDatasource{}
}

func (d *dao) NewApplicationDatasource() ApplicationDatasource {
	return &applicationDatasource{}
}

func (d *dao) NewUserAddressDatasource() UserAddressDatasource {
	return &userAddressDatasource{}
}

func (d *dao) NewOrderLifecycleDatasource() OrderLifecycleDatasource {
	return &orderLifecycleDatasource{}
}

func (d *dao) NewBusinessScheduleDatasource() BusinessScheduleDatasource {
	return &businessScheduleDatasource{}
}

func (d *dao) NewUnionOrderAndOrderedItemDatasource() UnionOrderAndOrderedItemDatasource {
	return &unionOrderAndOrderedItemDatasource{}
}

func (d *dao) NewUnionBusinessRoleAndUserDatasource() UnionBusinessRoleAndUserDatasource {
	return &unionBusinessRoleAndUserDatasource{}
}

func (d *dao) NewJwtTokenDatasource() JwtTokenDatasource {
	return &jwtTokenDatasource{}
}

func (d *dao) NewUnionBusinessRoleAndPermissionDatasource() UnionBusinessRoleAndPermissionDatasource {
	return &unionBusinessRoleAndPermissionDatasource{}
}

func (d *dao) NewBusinessAnalyticsDatasource() BusinessAnalyticsDatasource {
	return &businessAnalyticsDatasource{}
}

func (d *dao) NewItemAnalyticsDatasource() ItemAnalyticsDatasource {
	return &itemAnalyticsDatasource{}
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

func (d *dao) NewBusinessCollectionDatasource() BusinessCollectionDatasource {
	return &businessCollectionDatasource{}
}

func (d *dao) NewPartnerApplicationDatasource() PartnerApplicationDatasource {
	return &partnerApplicationDatasource{}
}

func (d *dao) NewBusinessRoleDatasource() BusinessRoleDatasource {
	return &businessRoleDatasource{}
}

func (d *dao) NewPermissionDatasource() PermissionDatasource {
	return &permissionDatasource{}
}

func (d *dao) NewUserPermissionDatasource() UserPermissionDatasource {
	return &userPermissionDatasource{}
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
