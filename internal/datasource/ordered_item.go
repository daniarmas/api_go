package datasource

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderedItemDatasource interface {
	BatchCreateOrderedItem(tx *gorm.DB, data *[]entity.OrderedItem) (*[]entity.OrderedItem, error)
	ListOrderedItemByIds(tx *gorm.DB, ids []uuid.UUID) (*[]entity.OrderedItem, error)
	ListOrderedItem(tx *gorm.DB, where *entity.OrderedItem) (*[]entity.OrderedItem, error)
}

type orderedItemDatasource struct{}

func (i *orderedItemDatasource) BatchCreateOrderedItem(tx *gorm.DB, data *[]entity.OrderedItem) (*[]entity.OrderedItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *orderedItemDatasource) ListOrderedItemByIds(tx *gorm.DB, ids []uuid.UUID) (*[]entity.OrderedItem, error) {
	var res []entity.OrderedItem
	result := tx.Where("id IN ? ", ids).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *orderedItemDatasource) ListOrderedItem(tx *gorm.DB, where *entity.OrderedItem) (*[]entity.OrderedItem, error) {
	var res []entity.OrderedItem
	result := tx.Model(&entity.OrderedItem{}).Where(where).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
