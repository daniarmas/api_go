package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionOrderAndOrderedItemDatasource interface {
	BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error)
	ListUnionOrderAndOrderedItem(tx *gorm.DB, where *models.UnionOrderAndOrderedItem, fields *[]string) (*[]models.UnionOrderAndOrderedItem, error)
}

type unionOrderAndOrderedItemDatasource struct{}

func (i *unionOrderAndOrderedItemDatasource) BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *unionOrderAndOrderedItemDatasource) ListUnionOrderAndOrderedItem(tx *gorm.DB, where *models.UnionOrderAndOrderedItem, fields *[]string) (*[]models.UnionOrderAndOrderedItem, error) {
	var res []models.UnionOrderAndOrderedItem
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
