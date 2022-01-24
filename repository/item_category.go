package repository

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ItemCategoryQuery interface {
	GetItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*models.BusinessItemCategory, error)
	ListItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*[]models.BusinessItemCategory, error)
	// SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error)
	// CreateItem(answer models.Item) (*int64, error)
	// UpdateItem(answer models.Item) (*models.Item, error)
	// DeleteItem(id int64) error
}

type itemCategoryQuery struct{}

func (i *itemCategoryQuery) ListItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*[]models.BusinessItemCategory, error) {
	var itemsCategory []models.BusinessItemCategory
	result := tx.Where(where).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}

func (i *itemCategoryQuery) GetItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*models.BusinessItemCategory, error) {
	var businessItemCategory *models.BusinessItemCategory
	result := tx.Where(where).Take(&businessItemCategory)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessItemCategory, nil
}
