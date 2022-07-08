package repository

import (
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	ListCartItem(tx *gorm.DB, where *entity.CartItem, cursor *time.Time, fields *[]string) (*[]entity.CartItem, error)
	ListCartItemAll(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*[]entity.CartItem, error)
	ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.CartItem, error)
	CreateCartItem(tx *gorm.DB, where *entity.CartItem) (*entity.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *entity.CartItem, data *entity.CartItem) (*entity.CartItem, error)
	GetCartItem(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*entity.CartItem, error)
	DeleteCartItem(tx *gorm.DB, where *entity.CartItem, ids *[]uuid.UUID) (*[]entity.CartItem, error)
	CartItemIsEmpty(tx *gorm.DB, where *entity.CartItem) (*bool, error)
}

type cartItemRepository struct{}

func (i *cartItemRepository) ListCartItem(tx *gorm.DB, where *entity.CartItem, cursor *time.Time, fields *[]string) (*[]entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItem(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) CartItemIsEmpty(tx *gorm.DB, where *entity.CartItem) (*bool, error) {
	res, err := Datasource.NewCartItemDatasource().CartItemIsEmpty(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) ListCartItemAll(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*[]entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItemAll(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItemInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) CreateCartItem(tx *gorm.DB, data *entity.CartItem) (*entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().CreateCartItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) UpdateCartItem(tx *gorm.DB, where *entity.CartItem, data *entity.CartItem) (*entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().UpdateCartItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) GetCartItem(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().GetCartItem(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *cartItemRepository) DeleteCartItem(tx *gorm.DB, where *entity.CartItem, ids *[]uuid.UUID) (*[]entity.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().DeleteCartItem(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
