package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessPaymentMethodTableName = "business_payment_method"

func (BusinessPaymentMethod) TableName() string {
	return BusinessPaymentMethodTableName
}

type BusinessPaymentMethod struct {
	ID              *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Type            string         `gorm:"column:type;not null"`
	Address         string         `gorm:"column:address"`
	BusinessId      *uuid.UUID     `gorm:"column:business_id;not null"`
	Business        Business       `gorm:"foreignKey:BusinessId"`
	PaymentMethodId *uuid.UUID     `gorm:"column:payment_method_id;not null"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodId"`
	CreateTime      time.Time      `gorm:"column:create_time;not null"`
	UpdateTime      time.Time      `gorm:"column:update_time;not null"`
	DeleteTime      gorm.DeletedAt `gorm:"index;column:delete_time"`
}

type BusinessPaymentMethodWithEnabled struct {
	ID              *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Type            string         `gorm:"column:type;not null"`
	Address         string         `gorm:"column:address"`
	Enabled         bool           `gorm:"column:enabled;not null"`
	BusinessId      *uuid.UUID     `gorm:"column:business_id;not null"`
	Business        Business       `gorm:"foreignKey:BusinessId"`
	PaymentMethodId *uuid.UUID     `gorm:"column:payment_method_id;not null"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodId"`
	CreateTime      time.Time      `gorm:"column:create_time;not null"`
	UpdateTime      time.Time      `gorm:"column:update_time;not null"`
	DeleteTime      gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *BusinessPaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
