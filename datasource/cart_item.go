package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartItemDatasource interface {
	ListCartItemAndItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItemAndItem, error)
	ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItem, error)
	CreateCartItem(tx *gorm.DB, where *models.CartItem) (*models.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error)
	ExistCartItem(tx *gorm.DB, where *models.CartItem) (*bool, error)
	GetCartItem(tx *gorm.DB, cartItem *models.CartItem) (*models.CartItem, error)
	DeleteCartItem(tx *gorm.DB, where *models.CartItem) error
}

type cartItemDatasource struct{}

func (i *cartItemDatasource) ListCartItemAndItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItemAndItem, error) {
	var cartItems []models.CartItemAndItem
	result := tx.Model(&models.CartItem{}).Limit(11).Select("cart_item.id, cart_item.name, cart_item.price, cart_item.quantity, cart_item.item_fk, cart_item.user_fk, cart_item.authorization_token_fk, item.thumbnail, item.thumbnail_blurhash, cart_item.create_time, cart_item.update_time").Joins("left join item on item.id = cart_item.item_fk").Where("cart_item.user_fk = ? AND cart_item.create_time < ?", where.UserFk, cursor).Order("cart_item.create_time desc").Scan(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cartItems, nil
}

func (i *cartItemDatasource) ListCartItem(tx *gorm.DB, where *models.CartItem, cursor *time.Time) (*[]models.CartItem, error) {
	var cartItems []models.CartItem
	result := tx.Limit(11).Where("cart_item.user_fk = ? AND cart_item.create_time > ?", where.UserFk, cursor).Order("cart_item.create_time desc").Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cartItems, nil
}

func (i *cartItemDatasource) CreateCartItem(tx *gorm.DB, data *models.CartItem) (*models.CartItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *cartItemDatasource) UpdateCartItem(tx *gorm.DB, where *models.CartItem, data *models.CartItem) (*models.CartItem, error) {
	result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}

func (v *cartItemDatasource) DeleteCartItem(tx *gorm.DB, where *models.CartItem) error {
	var cartItemResult *[]models.CartItem
	result := tx.Where(where).Delete(&cartItemResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *cartItemDatasource) ExistCartItem(tx *gorm.DB, where *models.CartItem) (*bool, error) {
	var existCartItemResult *models.CartItem
	var boolean = false
	result := tx.Where(where).Take(&existCartItemResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	if result.RowsAffected == 0 {
		return &boolean, nil
	} else {
		boolean = true
		return &boolean, nil
	}
}

func (v *cartItemDatasource) GetCartItem(tx *gorm.DB, cartItem *models.CartItem) (*models.CartItem, error) {
	var cartItemResult *models.CartItem
	result := tx.Where(cartItem).Take(&cartItemResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return cartItemResult, nil
}
