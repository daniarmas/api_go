package datastruct

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
	Email                    string         `gorm:"column:email"`
	Alias                    string         `gorm:"column:alias"`
	FullName                 string         `gorm:"column:fullname"`
	IsLegalAge               bool           `gorm:"column:is_legal_age"`
	HighQualityPhoto         string         `gorm:"column:high_quality_photo"`
	HighQualityPhotoBlurHash string         `gorm:"column:high_quality_photo_blurhash"`
	LowQualityPhoto          string         `gorm:"column:low_quality_photo"`
	LowQualityPhotoBlurHash  string         `gorm:"column:low_quality_photo_blurhash"`
	Thumbnail                string         `gorm:"column:thumbnail"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnail_blurhash"`
	CreateTime               time.Time      `gorm:"column:create_time"`
	UpdateTime               time.Time      `gorm:"column:update_time"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}
