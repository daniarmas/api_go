package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const OrderItemTableName = "ordered_item"

func (OrderedItem) TableName() string {
	return OrderItemTableName
}

type OrderedItem struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string         `gorm:"column:name;not null"`
	PriceCup   string         `gorm:"column:price_cup;not null"`
	CostCup    string         `gorm:"column:cost_cup"`
	ProfitCup  string         `gorm:"column:profit_cup"`
	PriceUsd   string         `gorm:"column:price_usd"`
	CostUsd    string         `gorm:"column:cost_usd"`
	ProfitUsd  string         `gorm:"column:profit_usd"`
	Quantity   int32          `gorm:"column:quantity;not null"`
	ItemId     *uuid.UUID     `gorm:"column:item_id;not null"`
	Item       CartItem       `gorm:"foreignKey:ItemId"`
	CartItemId *uuid.UUID     `gorm:"column:cart_item_id;not null"`
	CartItem   CartItem       `gorm:"foreignKey:CartItemId"`
	UserId     *uuid.UUID     `gorm:"column:user_id;not null"`
	User       User           `gorm:"foreignKey:UserId"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *OrderedItem) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *OrderedItem) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
