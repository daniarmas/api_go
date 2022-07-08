package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionDatasource interface {
	GetPermission(tx *gorm.DB, where *entity.Permission, fields *[]string) (*entity.Permission, error)
	ListPermissionAll(tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error)
	ListPermissionByIdAll(tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error)
}

type permissionDatasource struct{}

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

func (i *permissionDatasource) GetPermission(tx *gorm.DB, where *entity.Permission, fields *[]string) (*entity.Permission, error) {
	var res *entity.Permission
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
