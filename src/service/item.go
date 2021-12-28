package service

import (
	"github.com/daniarmas/api_go/src/datastruct"
	"github.com/daniarmas/api_go/src/repository"
)

type ItemService interface {
	GetItem(id string) (*datastruct.Item, error)
	ListItem() ([]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type itemService struct {
	dao repository.DAO
}

func NewItemService(dao repository.DAO) ItemService {
	return &itemService{dao: dao}
}

func (i *itemService) ListItem() ([]datastruct.Item, error) {
	items, err := i.dao.NewItemQuery().ListItem()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *itemService) GetItem(id string) (*datastruct.Item, error) {
	item, err := i.dao.NewItemQuery().GetItem(id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
