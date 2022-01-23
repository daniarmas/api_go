package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const DeviceTableName = "device"

func (Device) TableName() string {
	return DeviceTableName
}

type Device struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Platform                 string         `gorm:"column:platform;not null"`
	SystemVersion            string         `gorm:"column:system_version;not null"`
	DeviceId                 string         `gorm:"column:device_id;not null"`
	FirebaseCloudMessagingId string         `gorm:"column:firebase_cloud_messaging_id;not null"`
	Model                    string         `gorm:"column:model;not null"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *Device) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
