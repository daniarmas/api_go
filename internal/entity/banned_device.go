package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BannedDeviceTableName = "banned_device"

func (BannedDevice) TableName() string {
	return BannedDeviceTableName
}

type BannedDevice struct {
	ID                *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description       string         `gorm:"column:description;not null"`
	DeviceIdentifier  string         `gorm:"column:device_identifier;not null"`
	BanExpirationTime time.Time      `gorm:"column:ban_expiration_time;not null"`
	DeviceId          *uuid.UUID     `gorm:"column:device_id;not null"`
	Device            Device         `gorm:"foreignKey:DeviceId"`
	ModeratorId       *uuid.UUID     `gorm:"column:moderator_id;not null"`
	Moderator         User           `gorm:"foreignKey:ModeratorId"`
	CreateTime        time.Time      `gorm:"column:create_time;not null"`
	UpdateTime        time.Time      `gorm:"column:update_time;not null"`
	DeleteTime        gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BannedDevice) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}
