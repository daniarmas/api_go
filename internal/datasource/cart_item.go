package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartItemDatasource interface {
	ListCartItem(tx *gorm.DB, where *entity.CartItem, cursor *time.Time, fields *[]string) (*[]entity.CartItem, error)
	ListCartItemAll(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*[]entity.CartItem, error)
	ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.CartItem, error)
	CreateCartItem(tx *gorm.DB, where *entity.CartItem) (*entity.CartItem, error)
	UpdateCartItem(tx *gorm.DB, where *entity.CartItem, data *entity.CartItem) (*entity.CartItem, error)
	GetCartItem(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*entity.CartItem, error)
	DeleteCartItem(tx *gorm.DB, where *entity.CartItem, ids *[]uuid.UUID) (*[]entity.CartItem, error)
	CartItemIsEmpty(tx *gorm.DB, where *entity.CartItem) (*bool, error)
}

type cartItemDatasource struct{}

func (i *cartItemDatasource) ListCartItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.CartItem, error) {
	var cartItems []entity.CartItem
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where("id IN ?", ids).Select(*selectFields).Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cartItems, nil
}

func (i *cartItemDatasource) CartItemIsEmpty(tx *gorm.DB, where *entity.CartItem) (*bool, error) {
	var cartItems []entity.CartItem
	var res = true
	result := tx.Limit(1).Select("id").Where("cart_item.user_id = ?", where.UserId).Find(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(cartItems) == 0 {
		return &res, nil
	}
	res = false
	return &res, nil
}

func (i *cartItemDatasource) ListCartItem(tx *gorm.DB, where *entity.CartItem, cursor *time.Time, fields *[]string) (*[]entity.CartItem, error) {
	var res []entity.CartItem
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Model(&entity.CartItem{}).Select(*selectFields).Limit(11).Where("cart_item.user_id = ? AND cart_item.create_time < ?", where.UserId, cursor).Order("cart_item.create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *cartItemDatasource) ListCartItemAll(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*[]entity.CartItem, error) {
	var res []entity.CartItem
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where("cart_item.user_id = ?", where.UserId).Select(*selectFields).Order("cart_item.create_time desc").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *cartItemDatasource) CreateCartItem(tx *gorm.DB, data *entity.CartItem) (*entity.CartItem, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *cartItemDatasource) UpdateCartItem(tx *gorm.DB, where *entity.CartItem, data *entity.CartItem) (*entity.CartItem, error) {
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

func (v *cartItemDatasource) DeleteCartItem(tx *gorm.DB, where *entity.CartItem, ids *[]uuid.UUID) (*[]entity.CartItem, error) {
	var res *[]entity.CartItem
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *cartItemDatasource) GetCartItem(tx *gorm.DB, where *entity.CartItem, fields *[]string) (*entity.CartItem, error) {
	var res *entity.CartItem
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
