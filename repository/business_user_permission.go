package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessUserPermissionRepository interface {
	GetBusinessUserPermission(tx *gorm.DB, where *models.BusinessUserPermission, fields *[]string) (*models.BusinessUserPermission, error)
}

type businessUserPermissionRepository struct{}

func (v *businessUserPermissionRepository) GetBusinessUserPermission(tx *gorm.DB, where *models.BusinessUserPermission, fields *[]string) (*models.BusinessUserPermission, error) {
	res, err := Datasource.NewBusinessUserPermissionDatasource().GetBusinessUserPermission(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
