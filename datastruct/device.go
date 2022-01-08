package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const DeviceTableName = "Device"

func (Device) TableName() string {
	return DeviceTableName
}

type Device struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Platform                 string         `gorm:"column:platform"`
	SystemVersion            string         `gorm:"column:system_version"`
	DeviceId                 string         `gorm:"column:device_id"`
	FirebaseCloudMessagingId string         `gorm:"column:firebase_cloud_messaging_id"`
	Model                    string         `gorm:"column:model"`
	CreateTime               time.Time      `gorm:"column:create_time"`
	UpdateTime               time.Time      `gorm:"column:update_time"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}
