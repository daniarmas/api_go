package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessAnalyticsTableName = "business_analytics"

func (BusinessAnalytics) TableName() string {
	return OrderLifecycleTableName
}

type BusinessAnalytics struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Type       string         `gorm:"column:type"`
	BusinessId *uuid.UUID     `gorm:"column:business_id;not null"`
	Business   Business       `gorm:"foreignKey:BusinessId"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessAnalytics) BeforeCreate(tx *gorm.DB) (err error) {
	var nullTime time.Time
	if i.CreateTime == nullTime {
		i.CreateTime = time.Now().UTC()
	}
	if i.UpdateTime == nullTime {
		i.UpdateTime = time.Now().UTC()
	}
	return
}

func (u *BusinessAnalytics) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
