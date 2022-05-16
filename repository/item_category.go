package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessCollectionQuery interface {
	GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*models.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*[]models.BusinessCollection, error)
}

type businessCollectionQuery struct{}

func (i *businessCollectionQuery) ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*[]models.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().ListBusinessCollection(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessCollectionQuery) GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*models.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().GetBusinessCollection(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
