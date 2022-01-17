package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type ItemQuery interface {
	GetItem(id string) (datastruct.Item, error)
	ListItem(tx *gorm.DB) ([]datastruct.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type itemQuery struct{}

func (i *itemQuery) ListItem(tx *gorm.DB) ([]datastruct.Item, error) {
	var items []datastruct.Item
	result := tx.Limit(10).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (i *itemQuery) GetItem(id string) (datastruct.Item, error) {
	var item []datastruct.Item
	DB.Limit(1).Where("id = ?", id).Find(&item)
	if len(item) == 0 {
		return datastruct.Item{}, nil
	}
	return item[0], nil
}

func (i *itemQuery) SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]datastruct.Item, error) {
	var items []datastruct.Item
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
