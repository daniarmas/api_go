package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserPermissionRepository interface {
	GetUserPermission(tx *gorm.DB, where *models.UserPermission, fields *[]string) (*models.UserPermission, error)
}

type userPermissionRepository struct{}

func (v *userPermissionRepository) GetUserPermission(tx *gorm.DB, where *models.UserPermission, fields *[]string) (*models.UserPermission, error) {
	res, err := Datasource.NewUserPermissionDatasource().GetUserPermission(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
