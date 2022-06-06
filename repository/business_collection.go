package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessCollectionQuery interface {
	GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*models.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*[]models.BusinessCollection, error)
}

type businessCollectionQuery struct{}

func (i *businessCollectionQuery) ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*[]models.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().ListBusinessCollection(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessCollectionQuery) GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*models.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().GetBusinessCollection(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
