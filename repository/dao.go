package repository

import (
	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/utils"
	"gorm.io/gorm"
)

type DAO interface {
	NewItemQuery() ItemQuery
	NewVerificationCodeQuery() VerificationCodeQuery
	NewUserQuery() UserQuery
	NewBannedUserQuery() BannedUserQuery
	NewBannedDeviceQuery() BannedDeviceQuery
	NewDeviceQuery() DeviceQuery
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
}

type dao struct {
}

var Config *utils.Config
var Datasource datasource.DAO

func NewDAO(db *gorm.DB, config *utils.Config, datasourceDao datasource.DAO) DAO {
	Config = config
	Datasource = datasourceDao
	return &dao{}
}

// func NewConfig() (*utils.Config, error) {
// 	Config, err := utils.LoadConfig(".")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Config, nil
// }

// func NewDB(config *utils.Config) (*gorm.DB, error) {
// 	host := config.DBHost
// 	port := config.DBPort
// 	user := config.DBUser
// 	dbName := config.DBDatabase
// 	password := config.DBPassword

// 	// Starting a database
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
// 	newLogger := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags),
// 		logger.Config{
// 			SlowThreshold:             time.Millisecond * 200,
// 			LogLevel:                  logger.Info,
// 			IgnoreRecordNotFoundError: false,
// 			Colorful:                  true,
// 		},
// 	)
// 	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		SkipDefaultTransaction: true,
// 		Logger:                 newLogger,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return DB, nil
// }

func (d *dao) NewItemQuery() ItemQuery {
	return &itemQuery{}
}

func (d *dao) NewVerificationCodeQuery() VerificationCodeQuery {
	return &municipalityRepository{}
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

func (d *dao) NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (d *dao) NewBusinessScheduleRepository() BusinessScheduleRepository {
	return &businessScheduleRepository{}
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

func (d *dao) NewDeviceQuery() DeviceQuery {
	return &deviceQuery{}
}

func (d *dao) NewRefreshTokenQuery() RefreshTokenQuery {
	return &refreshTokenQuery{}
}

func (d *dao) NewAuthorizationTokenQuery() AuthorizationTokenQuery {
	return &authorizationTokenQuery{}
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
