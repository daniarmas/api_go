package repository

import (
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessRoleRepository interface {
	CreateBusinessRole(tx *gorm.DB, data *entity.BusinessRole) (*entity.BusinessRole, error)
	UpdateBusinessRole(tx *gorm.DB, where *entity.BusinessRole, data *entity.BusinessRole) (*entity.BusinessRole, error)
	GetBusinessRole(tx *gorm.DB, where *entity.BusinessRole, fields *[]string) (*entity.BusinessRole, error)
	ListBusinessRole(tx *gorm.DB, where *entity.BusinessRole, cursor *time.Time, fields *[]string) (*[]entity.BusinessRole, error)
	DeleteBusinessRole(tx *gorm.DB, where *entity.BusinessRole, ids *[]uuid.UUID) (*[]entity.BusinessRole, error)
}

type businessRoleRepository struct{}

func (v *businessRoleRepository) UpdateBusinessRole(tx *gorm.DB, where *entity.BusinessRole, data *entity.BusinessRole) (*entity.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().UpdateBusinessRole(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessRoleRepository) CreateBusinessRole(tx *gorm.DB, data *entity.BusinessRole) (*entity.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().CreateBusinessRole(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessRoleRepository) DeleteBusinessRole(tx *gorm.DB, where *entity.BusinessRole, ids *[]uuid.UUID) (*[]entity.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().DeleteBusinessRole(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessRoleRepository) ListBusinessRole(tx *gorm.DB, where *entity.BusinessRole, cursor *time.Time, fields *[]string) (*[]entity.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().ListBusinessRole(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessRoleRepository) GetBusinessRole(tx *gorm.DB, where *entity.BusinessRole, fields *[]string) (*entity.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().GetBusinessRole(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
