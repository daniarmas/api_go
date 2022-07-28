package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderedRepository interface {
	BatchCreateOrderedItem(tx *gorm.DB, data *[]entity.OrderedItem) (*[]entity.OrderedItem, error)
	ListOrderedItemByIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.OrderedItem, error)
	ListOrderedItem(tx *gorm.DB, where *entity.OrderedItem, fields *[]string) (*[]entity.OrderedItem, error)
}

type orderedRepository struct{}

func (i *orderedRepository) BatchCreateOrderedItem(tx *gorm.DB, data *[]entity.OrderedItem) (*[]entity.OrderedItem, error) {
	res, err := Datasource.NewOrderedItemDatasource().BatchCreateOrderedItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderedRepository) ListOrderedItemByIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.OrderedItem, error) {
	result, err := Datasource.NewOrderedItemDatasource().ListOrderedItemByIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *orderedRepository) ListOrderedItem(tx *gorm.DB, where *entity.OrderedItem, fields *[]string) (*[]entity.OrderedItem, error) {
	result, err := Datasource.NewOrderedItemDatasource().ListOrderedItem(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
