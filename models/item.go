package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ItemTableName = "item"

func (Item) TableName() string {
	return ItemTableName
}

type Item struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                     string         `gorm:"column:name;not null"`
	Description              string         `gorm:"column:description"`
	Price                    float64        `gorm:"column:price;not null"`
	Availability             int64          `gorm:"column:availability;not null"`
	BusinessFk               uuid.UUID      `gorm:"column:business_fk;not null"`
	Business                 Business       `gorm:"foreignKey:BusinessFk"`
	BusinessItemCategoryFk   uuid.UUID      `gorm:"column:business_item_category_fk;not null"`
	ProvinceFk               uuid.UUID      `gorm:"column:province_fk;not null"`
	MunicipalityFk           uuid.UUID      `gorm:"column:municipality_fk;not null"`
	HighQualityPhoto         string         `gorm:"column:high_quality_photo;not null"`
	HighQualityPhotoObject   string         `gorm:"column:high_quality_photo_object;not null"`
	HighQualityPhotoBlurHash string         `gorm:"column:high_quality_photo_blurhash;not null"`
	LowQualityPhoto          string         `gorm:"column:low_quality_photo;not null"`
	LowQualityPhotoObject    string         `gorm:"column:low_quality_photo_object;not null"`
	LowQualityPhotoBlurHash  string         `gorm:"column:low_quality_photo_blurhash;not null"`
	Thumbnail                string         `gorm:"column:thumbnail;not null"`
	ThumbnailObject          string         `gorm:"column:thumbnail_object;not null"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnail_blurhash;not null"`
	Cursor                   int32          `gorm:"column:cursor"`
	Status                   string         `gorm:"column:status;not null"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

type ItemBusiness struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                     string         `gorm:"column:name;not null"`
	Description              string         `gorm:"column:description"`
	Price                    float64        `gorm:"column:price;not null"`
	Availability             int64          `gorm:"column:availability;not null"`
	BusinessFk               uuid.UUID      `gorm:"column:business_fk;not null"`
	Business                 Business       `gorm:"foreignKey:BusinessFk"`
	BusinessItemCategoryFk   uuid.UUID      `gorm:"column:business_item_category_fk;not null"`
	HighQualityPhoto         string         `gorm:"column:high_quality_photo;not null"`
	HighQualityPhotoObject   string         `gorm:"column:high_quality_photo_object;not null"`
	HighQualityPhotoBlurHash string         `gorm:"column:high_quality_photo_blurhash;not null"`
	LowQualityPhoto          string         `gorm:"column:low_quality_photo;not null"`
	LowQualityPhotoObject    string         `gorm:"column:low_quality_photo_object;not null"`
	LowQualityPhotoBlurHash  string         `gorm:"column:low_quality_photo_blurhash;not null"`
	Thumbnail                string         `gorm:"column:thumbnail;not null"`
	ThumbnailObject          string         `gorm:"column:thumbnail_object;not null"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnail_blurhash;not null"`
	Cursor                   int32          `gorm:"column:cursor"`
	Status                   string         `gorm:"column:status;not null"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
	
}

// Contains tells whether a contains x.
func IndexById(items []Item, item Item) int {
	for i, n := range items {
		if n.ID == item.ID {
			return i
		}
	}
	return -1
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	var item Item
	result := tx.Select("cursor").Order("cursor desc").Last(&item)
	if result.Error != nil {
		return result.Error
	}
	i.Cursor = item.Cursor + 1
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
