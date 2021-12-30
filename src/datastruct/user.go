package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UserTableName = "user"

type User struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email                    string         `gorm:"column:email"`
	Alias                    string         `gorm:"column:alias"`
	FullName                 string         `gorm:"column:fullname"`
	IsLegalAge               bool           `gorm:"column:isLegalAge"`
	HighQualityPhoto         string         `gorm:"column:highQualityPhoto"`
	HighQualityPhotoBlurHash string         `gorm:"column:highQualityPhotoBlurHash"`
	LowQualityPhoto          string         `gorm:"column:lowQualityPhoto"`
	LowQualityPhotoBlurHash  string         `gorm:"column:lowQualityPhotoBlurHash"`
	Thumbnail                string         `gorm:"column:thumbnail"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnailBlurHash"`
	CreateTime               time.Time      `gorm:"column:createTime"`
	UpdateTime               time.Time      `gorm:"column:updateTime"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:deleteTime"`
}
