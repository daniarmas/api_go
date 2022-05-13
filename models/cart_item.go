package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const CartItemTableName = "cart_item"

func (CartItem) TableName() string {
	return CartItemTableName
}

type CartItem struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                 string         `gorm:"column:name;not null"`
	Price                float64        `gorm:"column:price;not null"`
	Quantity             int32          `gorm:"column:quantity;not null"`
	ItemId               uuid.UUID      `gorm:"column:item_id;not null"`
	BusinessId           uuid.UUID      `gorm:"column:business_id;not null"`
	Business             Business       `gorm:"foreignKey:BusinessId"`
	Thumbnail            string         `gorm:"column:thumbnail;not null"`
	ThumbnailBlurHash    string         `gorm:"column:thumbnail_blurhash;not null"`
	UserId               uuid.UUID      `gorm:"column:user_id;not null"`
	AuthorizationTokenId uuid.UUID      `gorm:"column:authorization_token_id;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *CartItem) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *CartItem) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
