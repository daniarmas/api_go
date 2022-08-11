package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type UnionOrderAndOrderedItemRepository interface {
	BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]entity.UnionOrderAndOrderedItem) (*[]entity.UnionOrderAndOrderedItem, error)
	ListUnionOrderAndOrderedItem(tx *gorm.DB, where *entity.UnionOrderAndOrderedItem) (*[]entity.UnionOrderAndOrderedItem, error)
}

type unionOrderAndOrderedItemRepository struct{}

func (i *unionOrderAndOrderedItemRepository) BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]entity.UnionOrderAndOrderedItem) (*[]entity.UnionOrderAndOrderedItem, error) {
	res, err := Datasource.NewUnionOrderAndOrderedItemDatasource().BatchCreateUnionOrderAndOrderedItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionOrderAndOrderedItemRepository) ListUnionOrderAndOrderedItem(tx *gorm.DB, where *entity.UnionOrderAndOrderedItem) (*[]entity.UnionOrderAndOrderedItem, error) {
	res, err := Datasource.NewUnionOrderAndOrderedItemDatasource().ListUnionOrderAndOrderedItem(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
