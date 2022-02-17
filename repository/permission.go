package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	PermissionExists(tx *gorm.DB, where *models.Permission) error
}

type permissionRepository struct{}

func (v *permissionRepository) PermissionExists(tx *gorm.DB, where *models.Permission) error {
	err := Datasource.NewPermissionDatasource().PermissionExists(tx, where)
	if err != nil {
		return err
	}
	return nil
}
