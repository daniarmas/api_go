package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetPermission(tx *gorm.DB, where *models.Permission, fields *[]string) (*models.Permission, error)
	ListPermissionAll(tx *gorm.DB, where *models.Permission) (*[]models.Permission, error)
	ListPermissionByIdAll(tx *gorm.DB, where *models.Permission, ids *[]uuid.UUID) (*[]models.Permission, error)
}

type permissionRepository struct{}

func (i *permissionRepository) GetPermission(tx *gorm.DB, where *models.Permission, fields *[]string) (*models.Permission, error) {
	res, err := Datasource.NewPermissionDatasource().GetPermission(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *permissionRepository) ListPermissionAll(tx *gorm.DB, where *models.Permission) (*[]models.Permission, error) {
	res, err := Datasource.NewPermissionDatasource().ListPermissionAll(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *permissionRepository) ListPermissionByIdAll(tx *gorm.DB, where *models.Permission, ids *[]uuid.UUID) (*[]models.Permission, error) {
	res, err := Datasource.NewPermissionDatasource().ListPermissionByIdAll(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
