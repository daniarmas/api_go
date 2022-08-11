package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessCollectionRepository interface {
	GetBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*entity.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*[]entity.BusinessCollection, error)
}

type businessCollectionRepository struct{}

func (i *businessCollectionRepository) ListBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*[]entity.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().ListBusinessCollection(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessCollectionRepository) GetBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*entity.BusinessCollection, error) {
	result, err := Datasource.NewBusinessCollectionDatasource().GetBusinessCollection(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
