package entity

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
	ID                    *uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                  string               `gorm:"column:name"`
	Address               string               `gorm:"column:address"`
	HighQualityPhoto      string               `gorm:"column:high_quality_photo;not null"`
	LowQualityPhoto       string               `gorm:"column:low_quality_photo;not null"`
	Thumbnail             string               `gorm:"column:thumbnail;not null"`
	BlurHash              string               `gorm:"column:blurhash;not null"`
	Cursor                int64                `gorm:"column:cursor"`
	Distance              float64              `gorm:"column:distance"`
	BusinessCategory      string               `gorm:"column:business_category"`
	TimeMarginOrderMonth  int32                `gorm:"column:time_margin_order_month"`
	TimeMarginOrderDay    int32                `gorm:"column:time_margin_order_day"`
	TimeMarginOrderHour   int32                `gorm:"column:time_margin_order_hour"`
	TimeMarginOrderMinute int32                `gorm:"column:time_margin_order_minute"`
	DeliveryPriceCup      string               `gorm:"column:delivery_price_cup"`
	ToPickUp              bool                 `gorm:"column:to_pick_up"`
	HomeDelivery          bool                 `gorm:"column:home_delivery"`
	Coordinates           ewkb.Point           `gorm:"column:coordinates"`
	ProvinceId            *uuid.UUID           `gorm:"column:province_id"`
	MunicipalityId        *uuid.UUID           `gorm:"column:municipality_id"`
	BusinessCategoryId    *uuid.UUID           `gorm:"column:business_category_id"`
	BusinessBrandId       *uuid.UUID           `gorm:"column:business_brand_id"`
	BusinessBrand         BusinessBrand        `gorm:"foreignKey:BusinessBrandId"`
	BusinessCollection    []BusinessCollection `gorm:"foreignKey:BusinessId"`
	BusinessSchedule      BusinessSchedule     `gorm:"foreignKey:BusinessId"`
	Status                string               `gorm:"column:status"`
	Municipality          []Municipality       `gorm:"many2many:union_business_and_municipality;"`
	CreateTime            time.Time            `gorm:"column:create_time"`
	UpdateTime            time.Time            `gorm:"column:update_time"`
	DeleteTime            gorm.DeletedAt       `gorm:"index;column:delete_time"`
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
