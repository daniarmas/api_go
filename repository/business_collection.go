package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessCollectionRepository interface {
	GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*models.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*[]models.BusinessCollection, error)
}

type businessCollectionRepository struct{}

func (i *businessCollectionRepository) ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*[]models.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().ListBusinessCollection(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessCollectionRepository) GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*models.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().GetBusinessCollection(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
