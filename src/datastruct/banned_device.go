package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BannedDeviceTableName = "BannedDevice"

type BannedDevice struct {
	ID                            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string         `gorm:"column:description"`
	DeviceId                      string         `gorm:"column:deviceId"`
	DeviceFk                      string         `gorm:"column:deviceFk"`
	ModeratorAuthorizationTokenFk string         `gorm:"column:moderatorAuthorizationTokenFk"`
	CreateTime                    time.Time      `gorm:"column:createTime"`
	UpdateTime                    time.Time      `gorm:"column:updateTime"`
	DeleteTime                    gorm.DeletedAt `gorm:"index;column:deleteTime"`
}
