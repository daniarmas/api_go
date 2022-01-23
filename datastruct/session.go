package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	RefreshTokenFk           uuid.UUID      `gorm:"type:uuid;column:refresh_token_fk;not null"`
	UserFk                   uuid.UUID      `gorm:"column:user_fk;not null"`
	User                     User           `gorm:"foreignKey:UserFk"`
	DeviceFk                 uuid.UUID      `gorm:"column:device_fk;not null"`
	Device                   Device         `gorm:"foreignKey:DeviceFk"`
	App                      string         `gorm:"column:app;not null"`
	AppVersion               string         `gorm:"column:app_version;not null"`
	Platform                 string         `gorm:"column:platform;not null"`
	SystemVersion            string         `gorm:"column:system_version;not null"`
	DeviceId                 string         `gorm:"column:device_id;not null"`
	FirebaseCloudMessagingId string         `gorm:"column:firebase_cloud_messaging_id;not null"`
	Model                    string         `gorm:"column:model;not null"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}
