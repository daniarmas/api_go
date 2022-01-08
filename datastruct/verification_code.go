package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const VerificationCodeTableName = "VerificationCode"

func (VerificationCode) TableName() string {
	return VerificationCodeTableName
}

type VerificationCode struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Code       string         `gorm:"column:code"`
	Email      string         `gorm:"column:email"`
	Type       string         `gorm:"column:type"`
	DeviceId   string         `gorm:"column:device_id"`
	CreateTime time.Time      `gorm:"column:create_time"`
	UpdateTime time.Time      `gorm:"column:update_time"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

// Note: Gorm will fail if the function signature
//  does not include `*gorm.DB` and `error`

func (vc *VerificationCode) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	vc.ID = uuid.New()
	return
}
