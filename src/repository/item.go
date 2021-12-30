package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type ItemQuery interface {
	GetItem(id string) (datastruct.Item, error)
	ListItem() ([]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type itemQuery struct{}

func (i *itemQuery) ListItem() ([]datastruct.Item, error) {
	var items []datastruct.Item
	DB.Table("Item").Limit(10).Find(&items)
	return items, nil
}

func (i *itemQuery) GetItem(id string) (datastruct.Item, error) {
	var item []datastruct.Item
	DB.Table("Item").Limit(1).Where("id = ?", id).Find(&item)
	if len(item) == 0 {
		return datastruct.Item{}, nil
	}
	return item[0], nil
}
