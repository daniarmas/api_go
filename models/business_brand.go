package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessBrandTableName = "business_brand"

func (BusinessBrand) TableName() string {
	return BusinessBrandTableName
}

type BusinessBrand struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name            string         `gorm:"column:name;not null"`
	BusinessOwnerId uuid.UUID      `gorm:"column:business_owner_id;not null"`
	CreateTime      time.Time      `gorm:"column:create_time;not null"`
	UpdateTime      time.Time      `gorm:"column:update_time;not null"`
	DeleteTime      gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *BusinessBrand) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *BusinessBrand) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
