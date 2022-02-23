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
	ItemFk               uuid.UUID      `gorm:"column:item_fk;not null"`
	BusinessFk           uuid.UUID      `gorm:"column:business_fk;not null"`
	Business             Business       `gorm:"foreignKey:BusinessFk"`
	UserFk               uuid.UUID      `gorm:"column:user_fk;not null"`
	AuthorizationTokenFk uuid.UUID      `gorm:"column:authorization_token_fk;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

type CartItemAndItem struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                 string         `gorm:"column:name;not null"`
	Price                float64        `gorm:"column:price;not null"`
	Quantity             int32          `gorm:"column:quantity;not null"`
	ItemFk               uuid.UUID      `gorm:"column:item_fk;not null"`
	UserFk               uuid.UUID      `gorm:"column:user_fk;not null"`
	AuthorizationTokenFk uuid.UUID      `gorm:"column:authorization_token_fk;not null"`
	Thumbnail            string         `gorm:"column:thumbnail;not null"`
	ThumbnailBlurHash    string         `gorm:"column:thumbnail_blurhash;not null"`
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
