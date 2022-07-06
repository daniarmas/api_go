package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionDatasource interface {
	GetPermission(tx *gorm.DB, where *models.Permission, fields *[]string) (*models.Permission, error)
	ListPermissionAll(tx *gorm.DB, where *models.Permission) (*[]models.Permission, error)
	ListPermissionByIdAll(tx *gorm.DB, where *models.Permission, ids *[]uuid.UUID) (*[]models.Permission, error)
}

type permissionDatasource struct{}

func (i *permissionDatasource) ListPermissionAll(tx *gorm.DB, where *models.Permission) (*[]models.Permission, error) {
	var res []models.Permission
	result := tx.Where(where).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *permissionDatasource) ListPermissionByIdAll(tx *gorm.DB, where *models.Permission, ids *[]uuid.UUID) (*[]models.Permission, error) {
	var res []models.Permission
	result := tx.Model(&models.Permission{}).Where(where).Where("id IN ?", *ids).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *permissionDatasource) GetPermission(tx *gorm.DB, where *models.Permission, fields *[]string) (*models.Permission, error) {
	var res *models.Permission
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
