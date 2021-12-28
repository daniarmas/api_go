package datastruct

const ItemTableName = "item"

type Item struct {
	ID                       string  `gorm:"column:id"`
	Name                     string  `gorm:"column:name"`
	Description              string  `gorm:"column:description"`
	Price                    float64  `gorm:"column:price"`
	Availability             int64   `gorm:"column:availability"`
	BusinessFk               string  `gorm:"column:businessFk"`
	BusinessItemCategoryFk   string  `gorm:"column:business"`
	HighQualityPhoto         string  `gorm:"column:highQualityPhoto"`
	HighQualityPhotoBlurHash string  `gorm:"column:highQualityPhotoBlurHash"`
	LowQualityPhoto          string  `gorm:"column:lowQualityPhoto"`
	LowQualityPhotoBlurHash  string  `gorm:"column:lowQualityPhotoBlurHash"`
	Thumbnail                string  `gorm:"column:thumbnail"`
	ThumbnailBlurHash        string  `gorm:"column:thumbnailBlurHash"`
	Cursor                   int64   `gorm:"column:cursor"`
}
