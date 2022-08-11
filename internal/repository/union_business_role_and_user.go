package repository

import (
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnionBusinessRoleAndUserRepository interface {
	ListUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, cursor *time.Time) (*[]entity.UnionBusinessRoleAndUser, error)
	ListUnionBusinessRoleAndUserInIds(tx *gorm.DB, ids []uuid.UUID) (*[]entity.UnionBusinessRoleAndUser, error)
	ListUnionBusinessRoleAndUserAll(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error)
	CreateUnionBusinessRoleAndUser(tx *gorm.DB, data *[]entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error)
	DeleteUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndUser, error)
	GetUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*entity.UnionBusinessRoleAndUser, error)
}

type unionBusinessRoleAndUserRepository struct{}

func (i *unionBusinessRoleAndUserRepository) ListUnionBusinessRoleAndUserAll(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().ListUnionBusinessRoleAndUserAll(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) ListUnionBusinessRoleAndUserInIds(tx *gorm.DB, ids []uuid.UUID) (*[]entity.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().ListUnionBusinessRoleAndUserInIds(tx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) ListUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, cursor *time.Time) (*[]entity.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().ListUnionBusinessRoleAndUser(tx, where, cursor)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) CreateUnionBusinessRoleAndUser(tx *gorm.DB, data *[]entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().CreateUnionBusinessRoleAndUser(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) DeleteUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().DeleteUnionBusinessRoleAndUser(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionBusinessRoleAndUserRepository) GetUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*entity.UnionBusinessRoleAndUser, error) {
	res, err := Datasource.NewUnionBusinessRoleAndUserDatasource().GetUnionBusinessRoleAndUser(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
