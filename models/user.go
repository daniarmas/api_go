package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UserTableName = "user"

func (User) TableName() string {
	return UserTableName
}

type User struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email                    string         `gorm:"column:email;not null"`
	Alias                    string         `gorm:"column:alias;not null"`
	FullName                 string         `gorm:"column:fullname;not null"`
	IsLegalAge               bool           `gorm:"column:is_legal_age;not null"`
	HighQualityPhoto         string         `gorm:"column:high_quality_photo"`
	HighQualityPhotoBlurHash string         `gorm:"column:high_quality_photo_blurhash"`
	LowQualityPhoto          string         `gorm:"column:low_quality_photo"`
	LowQualityPhotoBlurHash  string         `gorm:"column:low_quality_photo_blurhash"`
	Thumbnail                string         `gorm:"column:thumbnail"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnail_blurhash"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	return
}
