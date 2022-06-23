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
	ID                   *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                 string         `gorm:"column:name;not null"`
	Description          string         `gorm:"column:description"`
	PriceCup             string         `gorm:"column:price_cup;not null"`
	CostCup              string         `gorm:"column:cost_cup"`
	ProfitCup            string         `gorm:"column:profit_cup"`
	PriceUsd             string         `gorm:"column:price_usd"`
	CostUsd              string         `gorm:"column:cost_usd"`
	ProfitUsd            string         `gorm:"column:profit_usd"`
	AvailableFlag        bool           `gorm:"column:available_flag;not null"`
	EnabledFlag          bool           `gorm:"column:enabled_flag;not null"`
	Availability         int64          `gorm:"column:availability;not null"`
	BusinessId           *uuid.UUID     `gorm:"column:business_id;not null"`
	Business             Business       `gorm:"foreignKey:BusinessId"`
	BusinessCollectionId *uuid.UUID     `gorm:"column:business_collection_id;not null"`
	ProvinceId           *uuid.UUID     `gorm:"column:province_id;not null"`
	MunicipalityId       *uuid.UUID     `gorm:"column:municipality_id;not null"`
	HighQualityPhoto     string         `gorm:"column:high_quality_photo;not null"`
	LowQualityPhoto      string         `gorm:"column:low_quality_photo;not null"`
	Thumbnail            string         `gorm:"column:thumbnail;not null"`
	BlurHash             string         `gorm:"column:blurhash;not null"`
	Cursor               int32          `gorm:"column:cursor"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

type ItemBusiness struct {
	ID                   *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                 string         `gorm:"column:name;not null"`
	Description          string         `gorm:"column:description"`
	AvailableFlag        bool           `gorm:"column:available_flag;not null"`
	EnabledFlag          bool           `gorm:"column:enabled_flag;not null"`
	PriceCup             string         `gorm:"column:price_cup;not null"`
	CostCup              string         `gorm:"column:cost_cup;not null"`
	ProfitCup            string         `gorm:"column:profit_cup;not null"`
	PriceUsd             string         `gorm:"column:price_usd"`
	CostUsd              string         `gorm:"column:cost_usd"`
	ProfitUsd            string         `gorm:"column:profit_usd"`
	Availability         int64          `gorm:"column:availability;not null"`
	BusinessId           *uuid.UUID     `gorm:"column:business_id;not null"`
	Business             Business       `gorm:"foreignKey:BusinessId"`
	BusinessCollectionId *uuid.UUID     `gorm:"column:business_collection_id;not null"`
	HighQualityPhoto     string         `gorm:"column:high_quality_photo;not null"`
	LowQualityPhoto      string         `gorm:"column:low_quality_photo;not null"`
	Thumbnail            string         `gorm:"column:thumbnail;not null"`
	BlurHash             string         `gorm:"column:blurhash;not null"`
	Cursor               int32          `gorm:"column:cursor"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
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
