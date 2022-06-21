package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time, fields *[]string) (*[]models.CartItem, error)
	ListCartItemAll(tx *gorm.DB, where *models.CartItem, fields *[]string) (*[]models.CartItem, error)
	ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.CartItem, error)
	CreateCartItem(tx *gorm.DB, where *models.CartItem) (*models.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error)
	GetCartItem(tx *gorm.DB, where *models.CartItem, fields *[]string) (*models.CartItem, error)
	DeleteCartItem(tx *gorm.DB, where *models.CartItem, ids *[]uuid.UUID) (*[]models.CartItem, error)
	CartItemIsEmpty(tx *gorm.DB, where *models.CartItem) (*bool, error)
}

type cartItemRepository struct{}

func (i *cartItemRepository) ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time, fields *[]string) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItem(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) CartItemIsEmpty(tx *gorm.DB, where *models.CartItem) (*bool, error) {
	res, err := Datasource.NewCartItemDatasource().CartItemIsEmpty(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) ListCartItemAll(tx *gorm.DB, where *models.CartItem, fields *[]string) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItemAll(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().ListCartItemInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) CreateCartItem(tx *gorm.DB, data *models.CartItem) (*models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().CreateCartItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().UpdateCartItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *cartItemRepository) GetCartItem(tx *gorm.DB, where *models.CartItem, fields *[]string) (*models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().GetCartItem(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *cartItemRepository) DeleteCartItem(tx *gorm.DB, where *models.CartItem, ids *[]uuid.UUID) (*[]models.CartItem, error) {
	res, err := Datasource.NewCartItemDatasource().DeleteCartItem(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
