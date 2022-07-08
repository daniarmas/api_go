package entity

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
	ID                   *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                 string         `gorm:"column:name;not null"`
	PriceCup             string         `gorm:"column:price_cup;not null"`
	CostCup              string         `gorm:"column:cost_cup"`
	ProfitCup             string         `gorm:"column:profit_cup"`
	PriceUsd             string         `gorm:"column:price_usd"`
	CostUsd              string         `gorm:"column:cost_usd"`
	ProfitUsd             string         `gorm:"column:profit_usd"`
	Quantity             int32          `gorm:"column:quantity;not null"`
	ItemId               *uuid.UUID     `gorm:"column:item_id;not null"`
	BusinessId           *uuid.UUID     `gorm:"column:business_id;not null"`
	Business             Business       `gorm:"foreignKey:BusinessId"`
	Thumbnail            string         `gorm:"column:thumbnail;not null"`
	BlurHash             string         `gorm:"column:blurhash;not null"`
	UserId               *uuid.UUID     `gorm:"column:user_id;not null"`
	AuthorizationTokenId *uuid.UUID     `gorm:"column:authorization_token_id;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *CartItem) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}

func (u *CartItem) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
