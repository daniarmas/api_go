package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessUserRepository interface {
	GetBusinessUser(tx *gorm.DB, where *models.BusinessUser, fields *[]string) (*models.BusinessUser, error)
	CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error)
	DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser) error
}

type businessUserRepository struct{}

func (v *businessUserRepository) CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().CreateBusinessUser(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessUserRepository) GetBusinessUser(tx *gorm.DB, where *models.BusinessUser, fields *[]string) (*models.BusinessUser, error) {
	result, err := Datasource.NewBusinessUserDatasource().GetBusinessUser(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *businessUserRepository) DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser) error {
	err := Datasource.NewBusinessUserDatasource().DeleteBusinessUser(tx, where)
	if err != nil {
		return err
	}
	return nil
}
