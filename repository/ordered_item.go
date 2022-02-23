package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderedRepository interface {
	BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error)
	ListOrderedItemByIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.OrderedItem, error)
}

type orderedRepository struct{}

func (i *orderedRepository) BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error) {
	res, err := Datasource.NewOrderedItemDatasource().BatchCreateOrderedItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderedRepository) ListOrderedItemByIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.OrderedItem, error) {
	result, err := Datasource.NewOrderedItemDatasource().ListOrderedItemByIds(tx, ids)
	if err != nil {
		return nil, err
	}
	return result, nil
}
