package repository

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ItemQuery interface {
	GetItem(tx *gorm.DB, id string) (*models.Item, error)
	ListItem(tx *gorm.DB, where *models.Item) (*[]models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error)
	// CreateItem(answer models.Item) (*int64, error)
	// UpdateItem(answer models.Item) (*models.Item, error)
	// DeleteItem(id int64) error
}

type itemQuery struct{}

func (i *itemQuery) ListItem(tx *gorm.DB, where *models.Item) (*[]models.Item, error) {
	var items []models.Item
	result := tx.Limit(11).Where("business_fk = ? AND business_item_category_fk = ? AND cursor > ?", where.BusinessFk, where.BusinessItemCategoryFk, where.Cursor).Order("cursor asc").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (i *itemQuery) GetItem(tx *gorm.DB, id string) (*models.Item, error) {
	var item *models.Item
	result := DB.Preload("ItemPhoto").Where("item.id = ?", id).Take(&item)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return item, nil
}

func (i *itemQuery) SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error) {
	var items []models.Item
	var result *gorm.DB
	where := "%" + name + "%"
	if municipalityNotEqual {
		result = tx.Limit(int(limit+1)).Select("id, name, price, thumbnail, thumbnail_blurhash, cursor, status").Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_fk != ? AND province_fk = ? AND cursor > ?", where, municipalityFk, provinceFk, cursor).Order("cursor asc").Find(&items)
	} else {
		result = tx.Limit(int(limit+1)).Select("id, name, price, thumbnail, thumbnail_blurhash, cursor, status").Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_fk = ? AND province_fk = ? AND cursor > ?", where, municipalityFk, provinceFk, cursor).Order("cursor asc").Find(&items)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}
