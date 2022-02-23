package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderedRepository interface {
	BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error)
}

type orderedRepository struct{}

func (i *orderedRepository) BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error) {
	res, err := Datasource.NewOrderedItemDatasource().BatchCreateOrderedItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
