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
	var item datastruct.Item
	DB.Table("Item").Where("id = ?", id).First(&item)
	return item, nil
}
