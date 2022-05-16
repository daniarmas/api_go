package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserPermissionRepository interface {
	UserPermissionExists(tx *gorm.DB, where *models.UserPermission) error
}

type userPermissionRepository struct{}

func (v *userPermissionRepository) UserPermissionExists(tx *gorm.DB, where *models.UserPermission) error {
	err := Datasource.NewUserPermissionDatasource().UserPermissionExists(tx, where)
	if err != nil {
		return err
	}
	return nil
}
