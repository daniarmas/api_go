package datasource

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type UnionOrderAndOrderedItemDatasource interface {
	BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]entity.UnionOrderAndOrderedItem) (*[]entity.UnionOrderAndOrderedItem, error)
	ListUnionOrderAndOrderedItem(tx *gorm.DB, where *entity.UnionOrderAndOrderedItem, fields *[]string) (*[]entity.UnionOrderAndOrderedItem, error)
}

type unionOrderAndOrderedItemDatasource struct{}

func (i *unionOrderAndOrderedItemDatasource) BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]entity.UnionOrderAndOrderedItem) (*[]entity.UnionOrderAndOrderedItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *unionOrderAndOrderedItemDatasource) ListUnionOrderAndOrderedItem(tx *gorm.DB, where *entity.UnionOrderAndOrderedItem, fields *[]string) (*[]entity.UnionOrderAndOrderedItem, error) {
	var res []entity.UnionOrderAndOrderedItem
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