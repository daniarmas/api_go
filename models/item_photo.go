package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ItemPhotoTableName = "item_photo"

func (ItemPhoto) TableName() string {
	return ItemPhotoTableName
}

type ItemPhoto struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	ItemFk                   uuid.UUID      `gorm:"column:item_fk;not null"`
	HighQualityPhoto         string         `gorm:"column:high_quality_photo;not null"`
	HighQualityPhotoBlurHash string         `gorm:"column:high_quality_photo_blurhash;not null"`
	LowQualityPhoto          string         `gorm:"column:low_quality_photo;not null"`
	LowQualityPhotoBlurHash  string         `gorm:"column:low_quality_photo_blurhash;not null"`
	Thumbnail                string         `gorm:"column:thumbnail;not null"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnail_blurhash;not null"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *ItemPhoto) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
