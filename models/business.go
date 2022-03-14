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
	HighQualityPhoto         string                 `gorm:"column:high_quality_photo;not null"`
	HighQualityPhotoObject   string                 `gorm:"column:high_quality_photo_object;not null"`
	HighQualityPhotoBlurHash string                 `gorm:"column:high_quality_photo_blurhash;not null"`
	LowQualityPhoto          string                 `gorm:"column:low_quality_photo;not null"`
	LowQualityPhotoObject    string                 `gorm:"column:low_quality_photo_object;not null"`
	LowQualityPhotoBlurHash  string                 `gorm:"column:low_quality_photo_blurhash;not null"`
	Thumbnail                string                 `gorm:"column:thumbnail;not null"`
	ThumbnailObject          string                 `gorm:"column:thumbnail_object;not null"`
	ThumbnailBlurHash        string                 `gorm:"column:thumbnail_blurhash;not null"`
	Cursor                   int64                  `gorm:"column:cursor"`
	TimeMarginOrderMonth     int32                  `gorm:"column:time_margin_order_month"`
	TimeMarginOrderDay       int32                  `gorm:"column:time_margin_order_day"`
	TimeMarginOrderHour      int32                  `gorm:"column:time_margin_order_hour"`
	TimeMarginOrderMinute    int32                  `gorm:"column:time_margin_order_minute"`
	DeliveryPrice            float32                `gorm:"column:delivery_price"`
	ToPickUp                 bool                   `gorm:"column:to_pick_up"`
	HomeDelivery             bool                   `gorm:"column:home_delivery"`
	Coordinates              ewkb.Point             `gorm:"column:coordinates"`
	IsInRange                bool                   `gorm:"column:is_in_range"`
	ProvinceFk               uuid.UUID              `gorm:"column:province_fk"`
	MunicipalityFk           uuid.UUID              `gorm:"column:municipality_fk"`
	BusinessBrandFk          uuid.UUID              `gorm:"column:business_brand_fk"`
	BusinessItemCategory     []BusinessItemCategory `gorm:"foreignKey:BusinessFk"`
	BusinessSchedule         BusinessSchedule       `gorm:"foreignKey:BusinessFk"`
	Status                   string                 `gorm:"column:status"`
	Distance                 float64                `gorm:"column:distance"`
	Municipality             []Municipality         `gorm:"many2many:union_business_and_municipality;"`
	CreateTime               time.Time              `gorm:"column:create_time"`
	UpdateTime               time.Time              `gorm:"column:update_time"`
	DeleteTime               gorm.DeletedAt         `gorm:"index;column:delete_time"`
}

func (i *Business) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *Business) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
