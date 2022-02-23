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
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Price      float64        `gorm:"column:price;not null"`
	ItemFk     uuid.UUID      `gorm:"column:item_fk;not null"`
	Item       Item           `gorm:"foreignKey:ItemFk"`
	Quantity   int32          `gorm:"column:quantity;not null"`
	UserFk     uuid.UUID      `gorm:"column:user_fk;not null"`
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
