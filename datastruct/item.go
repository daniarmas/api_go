package datastruct

import "github.com/google/uuid"

const ItemTableName = "item"

type Item struct {
	ID                       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                     string    `gorm:"column:name"`
	Description              string    `gorm:"column:description"`
	Price                    float64   `gorm:"column:price"`
	Availability             int64     `gorm:"column:availability"`
	BusinessFk               uuid.UUID `gorm:"column:business_fk"`
	BusinessItemCategoryFk   uuid.UUID `gorm:"column:business"`
	HighQualityPhoto         string    `gorm:"column:high_quality_photo"`
	HighQualityPhotoBlurHash string    `gorm:"column:high_quality_photo_blurhash"`
	LowQualityPhoto          string    `gorm:"column:low_quality_photo"`
	LowQualityPhotoBlurHash  string    `gorm:"column:low_quality_photo_blurhash"`
	Thumbnail                string    `gorm:"column:thumbnail"`
	ThumbnailBlurHash        string    `gorm:"column:thumbnail_blurhash"`
	Cursor                   int64     `gorm:"column:cursor"`
}
