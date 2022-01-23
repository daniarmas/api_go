package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type ItemCategoryQuery interface {
	// GetItem(id string) (datastruct.Item, error)
	ListItemCategory(tx *gorm.DB, where *datastruct.BusinessItemCategory) (*[]datastruct.BusinessItemCategory, error)
	// SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type itemCategoryQuery struct{}

func (i *itemCategoryQuery) ListItemCategory(tx *gorm.DB, where *datastruct.BusinessItemCategory) (*[]datastruct.BusinessItemCategory, error) {
	var itemsCategory []datastruct.BusinessItemCategory
	result := tx.Where(where).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}
