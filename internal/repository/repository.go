package repository

import (
	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	NewUserConfigurationRepository() UserConfigurationRepository
	NewItemRepository() ItemRepository
	NewVerificationCodeRepository() VerificationCodeRepository
	NewUserRepository() UserRepository
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
	NewUserAddressRepository() UserAddressRepository
	NewPartnerApplicationRepository() PartnerApplicationRepository
	NewBusinessRoleRepository() BusinessRoleRepository
	NewJwtTokenRepository() JwtTokenRepository
	NewUnionBusinessRoleAndPermissionRepository() UnionBusinessRoleAndPermissionRepository
	NewUnionBusinessRoleAndUserRepository() UnionBusinessRoleAndUserRepository
	NewPermissionRepository() PermissionRepository
	NewPaymentMethodRepository() PaymentMethodRepository
	NewBusinessPaymentMethodRepository() BusinessPaymentMethodRepository
}

type repository struct {
}

var Datasource datasource.Datasource
var Rdb *redis.Client

func New(db *gorm.DB, config *config.Config, datasourceDao datasource.Datasource, rdb *redis.Client) Repository {
	Datasource = datasourceDao
	Rdb = rdb
	return &repository{}
}

func (d *repository) NewUserConfigurationRepository() UserConfigurationRepository {
	return &userConfigurationRepository{}
}

func (d *repository) NewJwtTokenRepository() JwtTokenRepository {
	return &jwtTokenRepository{}
}

func (d *repository) NewItemRepository() ItemRepository {
	return &itemRepository{}
}

func (d *repository) NewVerificationCodeRepository() VerificationCodeRepository {
	return &verificationCodeRepository{}
}

func (d *repository) NewApplicationRepository() ApplicationRepository {
	return &applicationRepository{}
}

func (d *repository) NewUserRepository() UserRepository {
	return &userRepository{}
}

func (d *repository) NewMunicipalityRepository() MunicipalityRepository {
	return &municipalityRepository{}
}

func (d *repository) NewBusinessUserRepository() BusinessUserRepository {
	return &businessUserRepository{}
}

func (d *repository) NewOrderLifecycleRepository() OrderLifecycleRepository {
	return &orderLifecycleRepository{}
}

func (d *repository) NewProvinceRepository() ProvinceRepository {
	return &provinceRepository{}
}

func (d *repository) NewUnionBusinessRoleAndPermissionRepository() UnionBusinessRoleAndPermissionRepository {
	return &unionBusinessRoleAndPermissionRepository{}
}
func (d *repository) NewUnionBusinessRoleAndUserRepository() UnionBusinessRoleAndUserRepository {
	return &unionBusinessRoleAndUserRepository{}
}

func (d *repository) NewPermissionRepository() PermissionRepository {
	return &permissionRepository{}
}

func (d *repository) NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (d *repository) NewBusinessScheduleRepository() BusinessScheduleRepository {
	return &businessScheduleRepository{}
}

func (d *repository) NewUserAddressRepository() UserAddressRepository {
	return &userAddressRepository{}
}

func (d *repository) NewOrderedRepository() OrderedRepository {
	return &orderedRepository{}
}

func (d *repository) NewUnionOrderAndOrderedItemRepository() UnionOrderAndOrderedItemRepository {
	return &unionOrderAndOrderedItemRepository{}
}

func (d *repository) NewDeviceRepository() DeviceRepository {
	return &deviceRepository{}
}

func (d *repository) NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{}
}

func (d *repository) NewPartnerApplicationRepository() PartnerApplicationRepository {
	return &partnerApplicationRepository{}
}

func (d *repository) NewAuthorizationTokenRepository() AuthorizationTokenRepository {
	return &authorizationTokenRepository{}
}

func (d *repository) NewBusinessCategoryRepository() BusinessCategoryRepository {
	return &businessCategoryRepository{}
}

func (d *repository) NewSessionRepository() SessionRepository {
	return &sessionRepository{}
}

func (d *repository) NewBusinessRepository() BusinessRepository {
	return &businessRepository{}
}

func (d *repository) NewBusinessCollectionRepository() BusinessCollectionRepository {
	return &businessCollectionRepository{}
}

func (d *repository) NewUnionBusinessAndMunicipalityRepository() UnionBusinessAndMunicipalityRepository {
	return &unionBusinessAndMunicipality{}
}

func (d *repository) NewBusinessRoleRepository() BusinessRoleRepository {
	return &businessRoleRepository{}
}

func (d *repository) NewObjectStorageRepository() ObjectStorageRepository {
	return &objectStorageRepository{}
}

func (d *repository) NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

func (d *repository) NewUserPermissionRepository() UserPermissionRepository {
	return &userPermissionRepository{}
}

func (d *repository) NewPaymentMethodRepository() PaymentMethodRepository {
	return &paymentMethodRepository{}
}

func (d *repository) NewBusinessPaymentMethodRepository() BusinessPaymentMethodRepository {
	return &businessPaymentMethodRepository{}
}
