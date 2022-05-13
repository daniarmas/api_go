package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const OrderLifecycleTableName = "order_lifecycle"

func (OrderLifecycle) TableName() string {
	return OrderLifecycleTableName
}

type OrderLifecycle struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Status     string         `gorm:"column:status"`
	OrderId    uuid.UUID      `gorm:"column:order_id;not null"`
	Order      Order          `gorm:"foreignKey:OrderId"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *OrderLifecycle) BeforeCreate(tx *gorm.DB) (err error) {
	var nullTime time.Time
	if i.CreateTime == nullTime {
		i.CreateTime = time.Now().UTC()
	}
	if i.UpdateTime == nullTime {
		i.UpdateTime = time.Now().UTC()
	}
	return
}

func (u *OrderLifecycle) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
