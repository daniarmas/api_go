package datasource

import (
	"github.com/google/uuid"

	"github.com/daniarmas/api_go/config"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type JsonWebTokenMetadata struct {
	TokenId *uuid.UUID
	Token   *string
}

type Datasource interface {
	NewUserConfigurationDatasource() UserConfigurationDatasource
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

type datasource struct {
	Gorm   *gorm.DB
	Config *config.Config
	Minio  *minio.Client
}

func New(db *gorm.DB, config *config.Config, minio *minio.Client) Datasource {
	return &datasource{
		Gorm:   db,
		Config: config,
		Minio:  minio,
	}
}

func (d *datasource) NewUserConfigurationDatasource() UserConfigurationDatasource {
	return &userConfigurationDatasource{}
}

func (d *datasource) NewObjectStorageDatasource() ObjectStorageDatasource {
	return &objectStorageDatasource{Minio: d.Minio}
}

func (d *datasource) NewBusinessCategoryDatasource() BusinessCategoryDatasource {
	return &businessCategoryDatasource{}
}

func (d *datasource) NewOrderedItemDatasource() OrderedItemDatasource {
	return &orderedItemDatasource{}
}

func (d *datasource) NewApplicationDatasource() ApplicationDatasource {
	return &applicationDatasource{}
}

func (d *datasource) NewUserAddressDatasource() UserAddressDatasource {
	return &userAddressDatasource{}
}

func (d *datasource) NewOrderLifecycleDatasource() OrderLifecycleDatasource {
	return &orderLifecycleDatasource{}
}

func (d *datasource) NewBusinessScheduleDatasource() BusinessScheduleDatasource {
	return &businessScheduleDatasource{}
}

func (d *datasource) NewUnionOrderAndOrderedItemDatasource() UnionOrderAndOrderedItemDatasource {
	return &unionOrderAndOrderedItemDatasource{}
}

func (d *datasource) NewUnionBusinessRoleAndUserDatasource() UnionBusinessRoleAndUserDatasource {
	return &unionBusinessRoleAndUserDatasource{}
}

func (d *datasource) NewJwtTokenDatasource() JwtTokenDatasource {
	return &jwtTokenDatasource{
		Config: d.Config,
	}
}

func (d *datasource) NewUnionBusinessRoleAndPermissionDatasource() UnionBusinessRoleAndPermissionDatasource {
	return &unionBusinessRoleAndPermissionDatasource{}
}

func (d *datasource) NewBusinessAnalyticsDatasource() BusinessAnalyticsDatasource {
	return &businessAnalyticsDatasource{}
}

func (d *datasource) NewItemAnalyticsDatasource() ItemAnalyticsDatasource {
	return &itemAnalyticsDatasource{}
}

func (d *datasource) NewBusinessUserDatasource() BusinessUserDatasource {
	return &businessUserDatasource{}
}

func (d *datasource) NewVerificationCodeDatasource() VerificationCodeDatasource {
	return &verificationCodeDatasource{}
}

func (d *datasource) NewUserDatasource() UserDatasource {
	return &userDatasource{}
}

func (d *datasource) NewRefreshTokenDatasource() RefreshTokenDatasource {
	return &refreshTokenDatasource{}
}

func (d *datasource) NewUnionBusinessAndMunicipalityDatasource() UnionBusinessAndMunicipalityDatasource {
	return &unionBusinessAndMunicipalityDatasource{}
}

func (d *datasource) NewAuthorizationTokenDatasource() AuthorizationTokenDatasource {
	return &authorizationTokenDatasource{}
}

func (d *datasource) NewItemDatasource() ItemDatasource {
	return &itemDatasource{}
}

func (d *datasource) NewProvinceDatasource() ProvinceDatasource {
	return &provinceDatasource{}
}

func (d *datasource) NewMunicipalityDatasource() MunicipalityDatasource {
	return &municipalityDatasource{}
}

func (d *datasource) NewBusinessCollectionDatasource() BusinessCollectionDatasource {
	return &businessCollectionDatasource{}
}

func (d *datasource) NewPartnerApplicationDatasource() PartnerApplicationDatasource {
	return &partnerApplicationDatasource{}
}

func (d *datasource) NewBusinessRoleDatasource() BusinessRoleDatasource {
	return &businessRoleDatasource{}
}

func (d *datasource) NewPermissionDatasource() PermissionDatasource {
	return &permissionDatasource{}
}

func (d *datasource) NewUserPermissionDatasource() UserPermissionDatasource {
	return &userPermissionDatasource{}
}

func (d *datasource) NewOrderDatasource() OrderDatasource {
	return &orderDatasource{}
}

func (d *datasource) NewDeviceDatasource() DeviceDatasource {
	return &deviceDatasource{}
}

func (d *datasource) NewBusinessDatasource() BusinessDatasource {
	return &businessDatasource{}
}

func (d *datasource) NewBannedUserDatasource() BannedUserDatasource {
	return &bannedUserDatasource{}
}

func (d *datasource) NewBannedDeviceDatasource() BannedDeviceDatasource {
	return &bannedDeviceDatasource{}
}

func (d *datasource) NewCartItemDatasource() CartItemDatasource {
	return &cartItemDatasource{}
}
