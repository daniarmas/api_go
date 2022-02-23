package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionOrderAndOrderedItemDatasource interface {
	BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error)
}

type unionOrderAndOrderedItemDatasource struct{}

func (i *unionOrderAndOrderedItemDatasource) BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
