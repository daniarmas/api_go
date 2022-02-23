package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemQuery interface {
	ListCartItemAndItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItemAndItem, error)
	ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItem, error)
	ListCartItemAll(tx *gorm.DB, where *models.CartItem) (*[]models.CartItem, error)
	ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID) (*[]models.CartItem, error)
	CreateCartItem(tx *gorm.DB, where *models.CartItem) (*models.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error)
	GetCartItem(tx *gorm.DB, cartItem *models.CartItem) (*models.CartItem, error)
	DeleteCartItem(tx *gorm.DB, where *models.CartItem) error
	CartItemQuantity(tx *gorm.DB, where *models.CartItem) (*bool, error)
}

type cartItemQuery struct{}

func (i *cartItemQuery) ListCartItemAndItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItemAndItem, error) {
	result, err := Datasource.NewCartItemDatasource().ListCartItemAndItem(tx, where, cursor)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItem, error) {
	result, err := Datasource.NewCartItemDatasource().ListCartItem(tx, where, cursor)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) CartItemQuantity(tx *gorm.DB, where *models.CartItem) (*bool, error) {
	result, err := Datasource.NewCartItemDatasource().CartItemQuantity(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) ListCartItemAll(tx *gorm.DB, where *models.CartItem) (*[]models.CartItem, error) {
	result, err := Datasource.NewCartItemDatasource().ListCartItemAll(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemQuery) ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID) (*[]models.CartItem, error) {
	result, err := Datasource.NewCartItemDatasource().ListCartItemInIds(tx, ids)
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

func (v *cartItemQuery) DeleteCartItem(tx *gorm.DB, where *models.CartItem) error {
	err := Datasource.NewCartItemDatasource().DeleteCartItem(tx, where)
	if err != nil {
		return err
	}
	return nil
}
