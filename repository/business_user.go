package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessUserRepository interface {
	GetBusinessUser(tx *gorm.DB, where *models.BusinessUser) (*models.BusinessUser, error)
	CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error)
	DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser, ids *[]uuid.UUID) (*[]models.BusinessUser, error)
}

type businessUserRepository struct{}

func (v *businessUserRepository) CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().CreateBusinessUser(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessUserRepository) GetBusinessUser(tx *gorm.DB, where *models.BusinessUser) (*models.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().GetBusinessUser(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessUserRepository) DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser, ids *[]uuid.UUID) (*[]models.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().DeleteBusinessUser(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
