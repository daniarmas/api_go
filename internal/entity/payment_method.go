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

func (r *PaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
