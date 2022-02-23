package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionOrderAndOrderedItemRepository interface {
	BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error)
	ListUnionOrderAndOrderedItem(tx *gorm.DB, where *models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error)
}

type unionOrderAndOrderedItemRepository struct{}

func (i *unionOrderAndOrderedItemRepository) BatchCreateUnionOrderAndOrderedItem(tx *gorm.DB, data *[]models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error) {
	res, err := Datasource.NewUnionOrderAndOrderedItemDatasource().BatchCreateUnionOrderAndOrderedItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *unionOrderAndOrderedItemRepository) ListUnionOrderAndOrderedItem(tx *gorm.DB, where *models.UnionOrderAndOrderedItem) (*[]models.UnionOrderAndOrderedItem, error) {
	res, err := Datasource.NewUnionOrderAndOrderedItemDatasource().ListUnionOrderAndOrderedItem(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
