package repository

import (

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnionBusinessRoleAndPermissionRepository interface {
	ListUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission) (*[]models.UnionBusinessRoleAndPermission, error)
	ListUnionBusinessRoleAndPermissionInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.UnionBusinessRoleAndPermission, error)
	DeleteUnionBusinessRoleAndPermissionByPermissionIds(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndPermission, error)
	CreateUnionBusinessRoleAndPermission(tx *gorm.DB, data *[]models.UnionBusinessRoleAndPermission) (*[]models.UnionBusinessRoleAndPermission, error)
	DeleteUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndPermission, error)
}

type unionBusinessRoleAndPermissionRepository struct{}

func (i *unionBusinessRoleAndPermissionRepository) ListUnionBusinessRoleAndPermissionInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.UnionBusinessRoleAndPermission, error) {
	res, err := Datasource.NewUnionBusinessRoleAndPermissionDatasource().ListUnionBusinessRoleAndPermissionInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndPermissionRepository) ListUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission) (*[]models.UnionBusinessRoleAndPermission, error) {
	res, err := Datasource.NewUnionBusinessRoleAndPermissionDatasource().ListUnionBusinessRoleAndPermission(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndPermissionRepository) CreateUnionBusinessRoleAndPermission(tx *gorm.DB, data *[]models.UnionBusinessRoleAndPermission) (*[]models.UnionBusinessRoleAndPermission, error) {
	res, err := Datasource.NewUnionBusinessRoleAndPermissionDatasource().CreateUnionBusinessRoleAndPermission(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndPermissionRepository) DeleteUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndPermission, error) {
	res, err := Datasource.NewUnionBusinessRoleAndPermissionDatasource().DeleteUnionBusinessRoleAndPermission(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndPermissionRepository) DeleteUnionBusinessRoleAndPermissionByPermissionIds(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndPermission, error) {
	res, err := Datasource.NewUnionBusinessRoleAndPermissionDatasource().DeleteUnionBusinessRoleAndPermissionByPermissionIds(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
