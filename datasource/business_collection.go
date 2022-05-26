package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessCollectionDatasource interface {
	GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*models.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*[]models.BusinessCollection, error)
}

type businessCollectionDatasource struct{}

func (i *businessCollectionDatasource) ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*[]models.BusinessCollection, error) {
	var itemsCategory []models.BusinessCollection
	result := tx.Where(where).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}

func (i *businessCollectionDatasource) GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection) (*models.BusinessCollection, error) {
	var businessBusinessCollection *models.BusinessCollection
	result := tx.Where(where).Take(&businessBusinessCollection)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessBusinessCollection, nil
}
