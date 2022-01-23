package models

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
	ID                            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string             `gorm:"column:description;not null"`
	DeviceId                      string             `gorm:"column:device_id;not null"`
	DeviceFk                      uuid.UUID          `gorm:"column:device_fk;not null"`
	Device                        Device             `gorm:"foreignKey:DeviceFk"`
	ModeratorAuthorizationTokenFk uuid.UUID          `gorm:"column:moderator_authorization_token_fk;not null"`
	AuthorizationToken            AuthorizationToken `gorm:"foreignKey:ModeratorAuthorizationTokenFk"`
	CreateTime                    time.Time          `gorm:"column:create_time;not null"`
	UpdateTime                    time.Time          `gorm:"column:update_time;not null"`
	DeleteTime                    gorm.DeletedAt     `gorm:"index;column:delete_time"`
}

func (i *BannedDevice) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
