package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const PaymentMethodTableName = "payment_method"

func (PaymentMethod) TableName() string {
	return PaymentMethodTableName
}

type PaymentMethod struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Type       string         `gorm:"column:type;not null"`
	Address    string         `gorm:"column:address"`
	Enabled    bool           `gorm:"column:enabled;not null"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *PaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	var nullTime time.Time
	if i.CreateTime == nullTime {
		i.CreateTime = time.Now().UTC()
	}
	if i.UpdateTime == nullTime {
		i.UpdateTime = time.Now().UTC()
	}
	return
}

func (u *PaymentMethod) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
