package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const BusinessTableName = "business"

func (Business) TableName() string {
	return BusinessTableName
}

type Business struct {
	ID                       uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                     string                 `gorm:"column:name"`
	Description              string                 `gorm:"column:description"`
	Address                  string                 `gorm:"column:address"`
	Phone                    string                 `gorm:"column:phone"`
	Email                    string                 `gorm:"column:email"`
	HighQualityPhoto         string                 `gorm:"column:high_quality_photo"`
	HighQualityPhotoBlurHash string                 `gorm:"column:high_quality_photo_blurhash"`
	LowQualityPhoto          string                 `gorm:"column:low_quality_photo"`
	LowQualityPhotoBlurHash  string                 `gorm:"column:low_quality_photo_blurhash"`
	Thumbnail                string                 `gorm:"column:thumbnail"`
	ThumbnailBlurHash        string                 `gorm:"column:thumbnail_blurhash"`
	Cursor                   int64                  `gorm:"column:cursor"`
	IsOpen                   bool                   `gorm:"column:is_open"`
	LeadDayTime              int32                  `gorm:"column:lead_day_time"`
	LeadHoursTime            int32                  `gorm:"column:lead_hours_time"`
	LeadMinutesTime          int32                  `gorm:"column:lead_minutes_time"`
	DeliveryPrice            float32                `gorm:"column:delivery_price"`
	ToPickUp                 bool                   `gorm:"column:to_pick_up"`
	HomeDelivery             bool                   `gorm:"column:home_delivery"`
	Coordinates              ewkb.Point             `gorm:"column:coordinates"`
	Polygon                  ewkb.Polygon           `gorm:"column:polygon"`
	IsInRange                bool                   `gorm:"column:is_in_range"`
	ProvinceFk               uuid.UUID              `gorm:"column:province_fk"`
	MunicipalityFk           uuid.UUID              `gorm:"column:municipality_fk"`
	BusinessBrandFk          uuid.UUID              `gorm:"column:business_brand_fk"`
	BusinessItemCategory     []BusinessItemCategory `gorm:"foreignKey:BusinessFk"`
	Status                   string                 `gorm:"column:status"`
	Distance                 float64                `gorm:"column:distance"`
	CreateTime               time.Time              `gorm:"column:create_time"`
	UpdateTime               time.Time              `gorm:"column:update_time"`
	DeleteTime               gorm.DeletedAt         `gorm:"index;column:delete_time"`
}

func (i *Business) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
