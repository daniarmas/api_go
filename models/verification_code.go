package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const VerificationCodeTableName = "verification_code"

func (VerificationCode) TableName() string {
	return VerificationCodeTableName
}

type VerificationCode struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Code       string         `gorm:"column:code;not null"`
	Email      string         `gorm:"column:email;not null"`
	Type       string         `gorm:"column:type;not null"`
	DeviceId   string         `gorm:"column:device_id;not null"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

// Note: Gorm will fail if the function signature
//  does not include `*gorm.DB` and `error`

func (vc *VerificationCode) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	vc.ID = uuid.New()
	vc.CreateTime = time.Now().UTC()
	vc.UpdateTime = time.Now().UTC()
	return
}

func (u *VerificationCode) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
