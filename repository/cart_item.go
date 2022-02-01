package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type CartItemQuery interface {
	ListCartItemAndItem(tx *gorm.DB, where *models.CartItem) (*[]models.CartItemAndItem, error)
}

type cartItemQuery struct{}

func (i *cartItemQuery) ListCartItemAndItem(tx *gorm.DB, where *models.CartItem) (*[]models.CartItemAndItem, error) {
	result, err := Datasource.NewCartItemDatasource().ListCartItemAndItem(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
