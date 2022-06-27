package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessRoleRepository interface {
	CreateBusinessRole(tx *gorm.DB, data *models.BusinessRole) (*models.BusinessRole, error)
	GetBusinessRole(tx *gorm.DB, where *models.BusinessRole, fields *[]string) (*models.BusinessRole, error)
	ListBusinessRole(tx *gorm.DB, where *models.BusinessRole, cursor *time.Time, fields *[]string) (*[]models.BusinessRole, error)
	DeleteBusinessRole(tx *gorm.DB, where *models.BusinessRole, ids *[]uuid.UUID) (*[]models.BusinessRole, error)
}

type businessRoleRepository struct{}

func (v *businessRoleRepository) CreateBusinessRole(tx *gorm.DB, data *models.BusinessRole) (*models.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().CreateBusinessRole(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessRoleRepository) DeleteBusinessRole(tx *gorm.DB, where *models.BusinessRole, ids *[]uuid.UUID) (*[]models.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().DeleteBusinessRole(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessRoleRepository) ListBusinessRole(tx *gorm.DB, where *models.BusinessRole, cursor *time.Time, fields *[]string) (*[]models.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().ListBusinessRole(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessRoleRepository) GetBusinessRole(tx *gorm.DB, where *models.BusinessRole, fields *[]string) (*models.BusinessRole, error) {
	res, err := Datasource.NewBusinessRoleDatasource().GetBusinessRole(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
