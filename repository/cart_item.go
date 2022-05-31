package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemQuery interface {
	ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItem, error)
	ListCartItemAll(tx *gorm.DB, where *models.CartItem) (*[]models.CartItem, error)
	ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.CartItem, error)
	CreateCartItem(tx *gorm.DB, where *models.CartItem) (*models.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error)
	GetCartItem(tx *gorm.DB, where *models.CartItem, fields *[]string) (*models.CartItem, error)
	DeleteCartItem(tx *gorm.DB, where *models.CartItem, ids *[]uuid.UUID) (*[]models.CartItem, error)
	CartItemQuantity(tx *gorm.DB, where *models.CartItem) (*bool, error)
}

type cartItemQuery struct{}

func (i *cartItemQuery) ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItem(tx, where, cursor)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemQuery) CartItemQuantity(tx *gorm.DB, where *models.CartItem) (*bool, error) {
	res, err := Datasource.NewCartItemDatasource().CartItemQuantity(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemQuery) ListCartItemAll(tx *gorm.DB, where *models.CartItem) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItemAll(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemQuery) ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItemInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemQuery) CreateCartItem(tx *gorm.DB, data *models.CartItem) (*models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().CreateCartItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemQuery) UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().UpdateCartItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemQuery) GetCartItem(tx *gorm.DB, where *models.CartItem, fields *[]string) (*models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().GetCartItem(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *cartItemQuery) DeleteCartItem(tx *gorm.DB, where *models.CartItem, ids *[]uuid.UUID) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().DeleteCartItem(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
