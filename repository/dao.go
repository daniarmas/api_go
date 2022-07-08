package repository

import (
	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/datasource"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type DAO interface {
	NewItemRepository() ItemRepository
	NewVerificationCodeRepository() VerificationCodeRepository
	NewUserRepository() UserRepository
	NewBannedUserRepository() BannedUserRepository
	NewBannedDeviceRepository() BannedDeviceRepository
	NewDeviceRepository() DeviceRepository
	NewRefreshTokenRepository() RefreshTokenRepository
	NewAuthorizationTokenRepository() AuthorizationTokenRepository
	NewSessionRepository() SessionRepository
	NewBusinessRepository() BusinessRepository
	NewMunicipalityRepository() MunicipalityRepository
	NewProvinceRepository() ProvinceRepository
	NewBusinessCollectionRepository() BusinessCollectionRepository
	NewObjectStorageRepository() ObjectStorageRepository
	NewCartItemRepository() CartItemRepository
	NewUserPermissionRepository() UserPermissionRepository
	NewUnionBusinessAndMunicipalityRepository() UnionBusinessAndMunicipalityRepository
	NewOrderRepository() OrderRepository
	NewOrderedRepository() OrderedRepository
	NewUnionOrderAndOrderedItemRepository() UnionOrderAndOrderedItemRepository
	NewBusinessUserRepository() BusinessUserRepository
	NewApplicationRepository() ApplicationRepository
	NewBusinessScheduleRepository() BusinessScheduleRepository
	NewOrderLifecycleRepository() OrderLifecycleRepository
	NewBusinessCategoryRepository() BusinessCategoryRepository
	NewBusinessAnalyticsRepository() BusinessAnalyticsRepository
	NewItemAnalyticsRepository() ItemAnalyticsRepository
	NewUserAddressRepository() UserAddressRepository
	NewPartnerApplicationRepository() PartnerApplicationRepository
	NewBusinessRoleRepository() BusinessRoleRepository
	NewJwtTokenRepository() JwtTokenRepository
	NewUnionBusinessRoleAndPermissionRepository() UnionBusinessRoleAndPermissionRepository
	NewUnionBusinessRoleAndUserRepository() UnionBusinessRoleAndUserRepository
	NewPermissionRepository() PermissionRepository
}

type dao struct {
}

var Config *config.Config
var Datasource datasource.DAO
var Rdb *redis.Client

func NewDAO(db *gorm.DB, config *config.Config, datasourceDao datasource.DAO, rdb *redis.Client) DAO {
	Rdb = rdb
	Config = config
	Datasource = datasourceDao
	return &dao{}
}

func (d *dao) NewJwtTokenRepository() JwtTokenRepository {
	return &jwtTokenRepository{}
}

func (d *dao) NewItemRepository() ItemRepository {
	return &itemRepository{}
}

func (d *dao) NewVerificationCodeRepository() VerificationCodeRepository {
	return &verificationCodeRepository{}
}

func (d *dao) NewApplicationRepository() ApplicationRepository {
	return &applicationRepository{}
}

func (d *dao) NewUserRepository() UserRepository {
	return &userRepository{}
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

func (d *dao) NewUnionBusinessRoleAndPermissionRepository() UnionBusinessRoleAndPermissionRepository {
	return &unionBusinessRoleAndPermissionRepository{}
}
func (d *dao) NewUnionBusinessRoleAndUserRepository() UnionBusinessRoleAndUserRepository {
	return &unionBusinessRoleAndUserRepository{}
}

func (d *dao) NewPermissionRepository() PermissionRepository {
	return &permissionRepository{}
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

func (d *dao) NewBannedUserRepository() BannedUserRepository {
	return &bannedUserRepository{}
}

func (d *dao) NewOrderedRepository() OrderedRepository {
	return &orderedRepository{}
}

func (d *dao) NewUnionOrderAndOrderedItemRepository() UnionOrderAndOrderedItemRepository {
	return &unionOrderAndOrderedItemRepository{}
}

func (d *dao) NewBannedDeviceRepository() BannedDeviceRepository {
	return &bannedDeviceRepository{}
}

func (d *dao) NewDeviceRepository() DeviceRepository {
	return &deviceRepository{}
}

func (d *dao) NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{}
}

func (d *dao) NewPartnerApplicationRepository() PartnerApplicationRepository {
	return &partnerApplicationRepository{}
}

func (d *dao) NewAuthorizationTokenRepository() AuthorizationTokenRepository {
	return &authorizationTokenRepository{}
}

func (d *dao) NewBusinessCategoryRepository() BusinessCategoryRepository {
	return &businessCategoryRepository{}
}

func (d *dao) NewSessionRepository() SessionRepository {
	return &sessionRepository{}
}

func (d *dao) NewBusinessRepository() BusinessRepository {
	return &businessRepository{}
}

func (d *dao) NewBusinessCollectionRepository() BusinessCollectionRepository {
	return &businessCollectionRepository{}
}

func (d *dao) NewUnionBusinessAndMunicipalityRepository() UnionBusinessAndMunicipalityRepository {
	return &unionBusinessAndMunicipality{}
}

func (d *dao) NewBusinessRoleRepository() BusinessRoleRepository {
	return &businessRoleRepository{}
}

func (d *dao) NewObjectStorageRepository() ObjectStorageRepository {
	return &objectStorageRepository{}
}

func (d *dao) NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

func (d *dao) NewUserPermissionRepository() UserPermissionRepository {
	return &userPermissionRepository{}
}
