package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type CartItemDatasource interface {
	ListCartItemAndItem(tx *gorm.DB, where *models.CartItem) (*[]models.CartItemAndItem, error)
}

type cartItemDatasource struct{}

func (i *cartItemDatasource) ListCartItemAndItem(tx *gorm.DB, where *models.CartItem) (*[]models.CartItemAndItem, error) {
	var cartItems []models.CartItemAndItem
	result := tx.Model(&models.CartItem{}).Limit(11).Select("cart_item.id, cart_item.name, cart_item.price, cart_item.quantity, cart_item.item_fk, cart_item.user_fk, cart_item.authorization_token_fk, item.thumbnail, item.thumbnail_blurhash, cart_item.cursor, cart_item.create_time, cart_item.update_time").Joins("left join item on item.id = cart_item.item_fk").Where("cart_item.user_fk = ? AND cart_item.cursor > ?", where.UserFk, where.Cursor).Order("cart_item.cursor asc").Scan(&cartItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cartItems, nil
}