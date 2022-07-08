package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetPermission(tx *gorm.DB, where *entity.Permission, fields *[]string) (*entity.Permission, error)
	ListPermissionAll(tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error)
	ListPermissionByIdAll(tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error)
}

type permissionRepository struct{}

func (i *permissionRepository) GetPermission(tx *gorm.DB, where *entity.Permission, fields *[]string) (*entity.Permission, error) {
	res, err := Datasource.NewPermissionDatasource().GetPermission(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *permissionRepository) ListPermissionAll(tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error) {
	res, err := Datasource.NewPermissionDatasource().ListPermissionAll(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *permissionRepository) ListPermissionByIdAll(tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error) {
	res, err := Datasource.NewPermissionDatasource().ListPermissionByIdAll(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
