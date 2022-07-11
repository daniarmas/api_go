package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UserConfigurationTableName = "user_configuration"

func (UserConfiguration) TableName() string {
	return UserConfigurationTableName
}

type UserConfiguration struct {
	ID                    *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	DataSaving            bool           `gorm:"column:data_saving"`
	HighQualityImagesWifi bool           `gorm:"column:high_quality_images_wifi"`
	HighQualityImagesData bool           `gorm:"column:high_quality_images_data"`
	PaymentMethod         string         `gorm:"column:payment_method"`
	UserId                *uuid.UUID     `gorm:"column:user_id;not null"`
	User                  User           `gorm:"foreignKey:UserId"`
	CreateTime            time.Time      `gorm:"column:create_time;not null"`
	UpdateTime            time.Time      `gorm:"column:update_time;not null"`
	DeleteTime            gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UserConfiguration) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UserConfiguration) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
