package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ItemAnalyticsTableName = "item_analytics"

func (ItemAnalytics) TableName() string {
	return OrderLifecycleTableName
}

type ItemAnalytics struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Type       string         `gorm:"column:type"`
	ItemId     *uuid.UUID     `gorm:"column:item_id;not null"`
	Item       Item           `gorm:"foreignKey:ItemId"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *ItemAnalytics) BeforeCreate(tx *gorm.DB) (err error) {
	var nullTime time.Time
	if i.CreateTime == nullTime {
		i.CreateTime = time.Now().UTC()
	}
	if i.UpdateTime == nullTime {
		i.UpdateTime = time.Now().UTC()
	}
	return
}

func (u *ItemAnalytics) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
