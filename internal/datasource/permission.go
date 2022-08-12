package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionDatasource interface {
	CreatePermission(tx *gorm.DB, data *entity.Permission) (*entity.Permission, error)
	ListPermission(tx *gorm.DB, where *entity.Permission, cursor time.Time) (*[]entity.Permission, error)
	GetPermission(tx *gorm.DB, where *entity.Permission) (*entity.Permission, error)
	ListPermissionAll(tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error)
	ListPermissionByIdAll(tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error)
}

type permissionDatasource struct{}

func (v *permissionDatasource) CreatePermission(tx *gorm.DB, data *entity.Permission) (*entity.Permission, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *permissionDatasource) ListPermission(tx *gorm.DB, where *entity.Permission, cursor time.Time) (*[]entity.Permission, error) {
	var res []entity.Permission
	result := tx.Model(&entity.Permission{}).Limit(11).Where(where).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *permissionDatasource) ListPermissionAll(tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error) {
	var res []entity.Permission
	result := tx.Where(where).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *permissionDatasource) ListPermissionByIdAll(tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error) {
	var res []entity.Permission
	result := tx.Model(&entity.Permission{}).Where(where).Where("id IN ?", *ids).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *permissionDatasource) GetPermission(tx *gorm.DB, where *entity.Permission) (*entity.Permission, error) {
	var res *entity.Permission
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
