package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const VerificationCodeTableName = "verification_code"

type VerificationCode struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Code       string         `gorm:"column:code"`
	Email      string         `gorm:"column:email"`
	Type       string         `gorm:"column:type"`
	DeviceId   string         `gorm:"column:deviceId"`
	CreateTime time.Time      `gorm:"column:createTime"`
	UpdateTime time.Time      `gorm:"column:updateTime"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:deleteTime"`
}

// Note: Gorm will fail if the function signature
//  does not include `*gorm.DB` and `error`

func (vc *VerificationCode) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	vc.ID = uuid.New()
	return
}
