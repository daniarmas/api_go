package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TermsViolatedBannedDeviceTableName = "terms_violated_banned_device"

func (TermsViolatedBannedDevice) TableName() string {
	return TermsViolatedBannedDeviceTableName
}

type TermsViolatedBannedDevice struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	BannedDeviceId uuid.UUID      `gorm:"column:banned_device_id;not null"`
	BannedDevice   BannedDevice   `gorm:"foreignKey:BannedDeviceId"`
	TermId         uuid.UUID      `gorm:"column:term_id;not null"`
	Term           Term           `gorm:"foreignKey:TermId"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *TermsViolatedBannedDevice) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *TermsViolatedBannedDevice) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
