package repository

import (
	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/utils"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type DAO interface {
	NewItemQuery() ItemQuery
	NewVerificationCodeQuery() VerificationCodeRepository
	NewUserQuery() UserQuery
	NewBannedUserQuery() BannedUserQuery
	NewBannedDeviceQuery() BannedDeviceQuery
	NewDeviceQuery() DeviceRepository
	NewRefreshTokenQuery() RefreshTokenQuery
	NewAuthorizationTokenQuery() AuthorizationTokenQuery
	NewSessionQuery() SessionQuery
	NewBusinessQuery() BusinessQuery
	NewMunicipalityRepository() MunicipalityRepository
	NewProvinceRepository() ProvinceRepository
	NewBusinessCollectionQuery() BusinessCollectionQuery
	NewObjectStorageRepository() ObjectStorageRepository
	NewCartItemRepository() CartItemQuery
	NewUserPermissionRepository() UserPermissionRepository
	NewUnionBusinessAndMunicipalityRepository() UnionBusinessAndMunicipalityRepository
	NewOrderRepository() OrderRepository
	NewOrderedRepository() OrderedRepository
	NewUnionOrderAndOrderedItemRepository() UnionOrderAndOrderedItemRepository
	NewBusinessUserRepository() BusinessUserRepository
	NewDeprecatedVersionAppRepository() DeprecatedVersionAppRepository
	NewBusinessScheduleRepository() BusinessScheduleRepository
	NewOrderLifecycleRepository() OrderLifecycleRepository
	NewBusinessCategoryRepository() BusinessCategoryRepository
	NewBusinessAnalyticsRepository() BusinessAnalyticsRepository
	NewItemAnalyticsRepository() ItemAnalyticsRepository
	NewUserAddressRepository() UserAddressRepository
}

type dao struct {
}

var Config *utils.Config
var Datasource datasource.DAO
var Rdb *redis.Client

func NewDAO(db *gorm.DB, config *utils.Config, datasourceDao datasource.DAO, rdb *redis.Client) DAO {
	Rdb = rdb
	Config = config
	Datasource = datasourceDao
	return &dao{}
}

func (d *dao) NewItemQuery() ItemQuery {
	return &itemQuery{}
}

func (d *dao) NewVerificationCodeQuery() VerificationCodeRepository {
	return &verificationCodeRepository{}
}

func (d *dao) NewDeprecatedVersionAppRepository() DeprecatedVersionAppRepository {
	return &deprecatedVersionAppRepository{}
}

func (d *dao) NewUserQuery() UserQuery {
	return &userQuery{}
}

func (d *dao) NewMunicipalityRepository() MunicipalityRepository {
	return &municipalityRepository{}
}

func (d *dao) NewBusinessUserRepository() BusinessUserRepository {
	return &businessUserRepository{}
}

func (d *dao) NewOrderLifecycleRepository() OrderLifecycleRepository {
	return &orderLifecycleRepository{}
}

func (d *dao) NewProvinceRepository() ProvinceRepository {
	return &provinceRepository{}
}

func (d *dao) NewBusinessAnalyticsRepository() BusinessAnalyticsRepository {
	return &businessAnalyticsRepository{}
}

func (d *dao) NewItemAnalyticsRepository() ItemAnalyticsRepository {
	return &itemAnalyticsRepository{}
}

func (d *dao) NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (d *dao) NewBusinessScheduleRepository() BusinessScheduleRepository {
	return &businessScheduleRepository{}
}

func (d *dao) NewUserAddressRepository() UserAddressRepository {
	return &userAddressRepository{}
}

func (d *dao) NewBannedUserQuery() BannedUserQuery {
	return &bannedUserQuery{}
}

func (d *dao) NewOrderedRepository() OrderedRepository {
	return &orderedRepository{}
}

func (d *dao) NewUnionOrderAndOrderedItemRepository() UnionOrderAndOrderedItemRepository {
	return &unionOrderAndOrderedItemRepository{}
}

func (d *dao) NewBannedDeviceQuery() BannedDeviceQuery {
	return &bannedDeviceQuery{}
}

func (d *dao) NewDeviceQuery() DeviceRepository {
	return &deviceRepository{}
}

func (d *dao) NewRefreshTokenQuery() RefreshTokenQuery {
	return &refreshTokenQuery{}
}

func (d *dao) NewAuthorizationTokenQuery() AuthorizationTokenQuery {
	return &authorizationTokenQuery{}
}

func (d *dao) NewBusinessCategoryRepository() BusinessCategoryRepository {
	return &businessCategoryRepository{}
}

func (d *dao) NewSessionQuery() SessionQuery {
	return &sessionQuery{}
}

func (d *dao) NewBusinessQuery() BusinessQuery {
	return &businessQuery{}
}

func (d *dao) NewBusinessCollectionQuery() BusinessCollectionQuery {
	return &businessCollectionQuery{}
}

func (d *dao) NewUnionBusinessAndMunicipalityRepository() UnionBusinessAndMunicipalityRepository {
	return &unionBusinessAndMunicipality{}
}

func (d *dao) NewObjectStorageRepository() ObjectStorageRepository {
	return &objectStorageRepository{}
}

func (d *dao) NewCartItemRepository() CartItemQuery {
	return &cartItemQuery{}
}

func (d *dao) NewUserPermissionRepository() UserPermissionRepository {
	return &userPermissionRepository{}
}
