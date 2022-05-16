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
	ID                       uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email                    string           `gorm:"column:email;not null"`
	PhoneNumber              string           `gorm:"column:phone_number"`
	UserAddress              []UserAddress    `gorm:"foreignKey:UserFk"`
	UserPermissions          []UserPermission `gorm:"foreignKey:UserFk"`
	FullName                 string           `gorm:"column:fullname;not null"`
	IsLegalAge               bool             `gorm:"column:is_legal_age;not null"`
	HighQualityPhoto         string           `gorm:"column:high_quality_photo;not null"`
	HighQualityPhotoBlurHash string           `gorm:"column:high_quality_photo_blurhash;not null"`
	LowQualityPhoto          string           `gorm:"column:low_quality_photo;not null"`
	LowQualityPhotoBlurHash  string           `gorm:"column:low_quality_photo_blurhash;not null"`
	Thumbnail                string           `gorm:"column:thumbnail;not null"`
	ThumbnailBlurHash        string           `gorm:"column:thumbnail_blurhash;not null"`
	CreateTime               time.Time        `gorm:"column:create_time;not null"`
	UpdateTime               time.Time        `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt   `gorm:"index;column:delete_time"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
