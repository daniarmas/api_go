package datasource

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderedItemDatasource interface {
	BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error)
	ListOrderedItemByIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.OrderedItem, error)
	ListOrderedItem(tx *gorm.DB, where *models.OrderedItem) (*[]models.OrderedItem, error)
}

type orderedItemDatasource struct{}

func (i *orderedItemDatasource) BatchCreateOrderedItem(tx *gorm.DB, data *[]models.OrderedItem) (*[]models.OrderedItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *orderedItemDatasource) ListOrderedItemByIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.OrderedItem, error) {
	var response []models.OrderedItem
	result := tx.Where("id IN ? ", *ids).Find(&response)
	if result.Error != nil {
		return nil, result.Error
	}
	return &response, nil
}

func (i *orderedItemDatasource) ListOrderedItem(tx *gorm.DB, where *models.OrderedItem) (*[]models.OrderedItem, error) {
	var response []models.OrderedItem
	result := tx.Model(&models.OrderedItem{}).Where(where).Find(&response)
	if result.Error != nil {
		return nil, result.Error
	}
	return &response, nil
}
