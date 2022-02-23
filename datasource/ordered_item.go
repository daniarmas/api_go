package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderedItemDatasource interface {
	BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error)
}

type orderedItemDatasource struct{}

func (i *orderedItemDatasource) BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
