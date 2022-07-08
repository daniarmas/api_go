package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const RefreshTokenTableName = "refresh_token"

func (RefreshToken) TableName() string {
	return RefreshTokenTableName
}

type RefreshToken struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId     *uuid.UUID     `gorm:"column:user_id;not null"`
	User       User           `gorm:"foreignKey:UserId"`
	DeviceId   *uuid.UUID     `gorm:"column:device_id;not null"`
	Device     Device         `gorm:"foreignKey:DeviceId"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
