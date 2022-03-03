package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessUserTableName = "business_user"

func (BusinessUser) TableName() string {
	return BusinessUserTableName
}

type BusinessUser struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	IsBusinessOwner bool           `gorm:"column:is_business_owner;not null"`
	UserFk          uuid.UUID      `gorm:"column:user_fk;not null"`
	CreateTime      time.Time      `gorm:"column:create_time;not null"`
	UpdateTime      time.Time      `gorm:"column:update_time;not null"`
	DeleteTime      gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *BusinessUser) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *BusinessUser) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
