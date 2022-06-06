package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessCollectionDatasource interface {
	GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*models.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*[]models.BusinessCollection, error)
}

type businessCollectionDatasource struct{}

func (i *businessCollectionDatasource) ListBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*[]models.BusinessCollection, error) {
	var itemsCategory []models.BusinessCollection
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}

func (i *businessCollectionDatasource) GetBusinessCollection(tx *gorm.DB, where *models.BusinessCollection, fields *[]string) (*models.BusinessCollection, error) {
	var businessBusinessCollection *models.BusinessCollection
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&businessBusinessCollection)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessBusinessCollection, nil
}
