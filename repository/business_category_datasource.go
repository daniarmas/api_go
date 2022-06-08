package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessCategoryRepository interface {
	GetBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, fields *[]string) (*models.BusinessCategory, error)
	CreateBusinessCategory(tx *gorm.DB, data *models.BusinessCategory) (*models.BusinessCategory, error)
	DeleteBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, ids *[]uuid.UUID) (*[]models.BusinessCategory, error)
}

type businessCategoryRepository struct{}

func (v *businessCategoryRepository) CreateBusinessCategory(tx *gorm.DB, data *models.BusinessCategory) (*models.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().CreateBusinessCategory(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessCategoryRepository) GetBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, fields *[]string) (*models.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().GetBusinessCategory(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessCategoryRepository) DeleteBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, ids *[]uuid.UUID) (*[]models.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().DeleteBusinessCategory(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
