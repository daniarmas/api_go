package dto

type Item struct {
	ID                       string
	Name                     string
	Description              string
	Price                    float64
	Availability             int64
	BusinessFk               string
	BusinessItemCategoryFk   string
	HighQualityPhoto         string
	HighQualityPhotoBlurHash string
	LowQualityPhoto          string
	LowQualityPhotoBlurHash  string
	Thumbnail                string
	ThumbnailBlurHash        string
	Cursor                   int64
	Photos                   []ItemPhoto
	// Status ItemStatusType
}

type ItemPhoto struct {
	Id                       string
	ItemFk                   string
	HighQualityPhoto         string
	HighQualityPhotoBlurHash string
	LowQualityPhoto          string
	LowQualityPhotoBlurHash  string
	Thumbnail                string
	ThumbnailBlurHash        string
}
