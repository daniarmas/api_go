package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type CartItemQuery interface {
	ListCartItemAndItem(tx *gorm.DB, where *models.CartItem) (*[]models.CartItemAndItem, error)
	CreateCartItem(tx *gorm.DB, where *models.CartItem) (*models.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error)
	GetCartItem(tx *gorm.DB, cartItem *models.CartItem) (*models.CartItem, error)
}

type cartItemQuery struct{}

func (i *cartItemQuery) ListCartItemAndItem(tx *gorm.DB, where *models.CartItem) (*[]models.CartItemAndItem, error) {
	result, err := Datasource.NewCartItemDatasource().ListCartItemAndItem(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) CreateCartItem(tx *gorm.DB, data *models.CartItem) (*models.CartItem, error) {
	result, err := Datasource.NewCartItemDatasource().CreateCartItem(tx, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error) {
	result, err := Datasource.NewCartItemDatasource().UpdateCartItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) GetCartItem(tx *gorm.DB, where *models.CartItem) (*models.CartItem, error) {
	result, err := Datasource.NewCartItemDatasource().GetCartItem(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
