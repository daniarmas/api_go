package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnionBusinessRoleAndUserRepository interface {
	ListUnionBusinessRoleAndUser(tx *gorm.DB, where *models.UnionBusinessRoleAndUser, cursor *time.Time, fields *[]string) (*[]models.UnionBusinessRoleAndUser, error)
	ListUnionBusinessRoleAndUserInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.UnionBusinessRoleAndUser, error)
	ListUnionBusinessRoleAndUserAll(tx *gorm.DB, where *models.UnionBusinessRoleAndUser) (*[]models.UnionBusinessRoleAndUser, error)
	CreateUnionBusinessRoleAndUser(tx *gorm.DB, data *[]models.UnionBusinessRoleAndUser) (*[]models.UnionBusinessRoleAndUser, error)
	DeleteUnionBusinessRoleAndUser(tx *gorm.DB, where *models.UnionBusinessRoleAndUser, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndUser, error)
	GetUnionBusinessRoleAndUser(tx *gorm.DB, where *models.UnionBusinessRoleAndUser) (*models.UnionBusinessRoleAndUser, error)
}

type unionBusinessRoleAndUserRepository struct{}

func (i *unionBusinessRoleAndUserRepository) ListUnionBusinessRoleAndUserAll(tx *gorm.DB, where *models.UnionBusinessRoleAndUser) (*[]models.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().ListUnionBusinessRoleAndUserAll(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) ListUnionBusinessRoleAndUserInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().ListUnionBusinessRoleAndUserInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) ListUnionBusinessRoleAndUser(tx *gorm.DB, where *models.UnionBusinessRoleAndUser, cursor *time.Time, fields *[]string) (*[]models.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().ListUnionBusinessRoleAndUser(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) CreateUnionBusinessRoleAndUser(tx *gorm.DB, data *[]models.UnionBusinessRoleAndUser) (*[]models.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().CreateUnionBusinessRoleAndUser(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) DeleteUnionBusinessRoleAndUser(tx *gorm.DB, where *models.UnionBusinessRoleAndUser, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().DeleteUnionBusinessRoleAndUser(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) GetUnionBusinessRoleAndUser(tx *gorm.DB, where *models.UnionBusinessRoleAndUser) (*models.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().GetUnionBusinessRoleAndUser(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
