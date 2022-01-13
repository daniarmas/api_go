package datastruct

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
	ID                            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string         `gorm:"column:description"`
	DeviceId                      string         `gorm:"column:device_id"`
	DeviceFk                      uuid.UUID      `gorm:"column:device_fk"`
	ModeratorAuthorizationTokenFk uuid.UUID      `gorm:"column:moderator_authorization_token_fk"`
	CreateTime                    time.Time      `gorm:"column:create_time"`
	UpdateTime                    time.Time      `gorm:"column:update_time"`
	DeleteTime                    gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BannedDevice) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
