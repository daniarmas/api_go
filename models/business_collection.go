package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessCollectionTableName = "business_collection"

func (BusinessCollection) TableName() string {
	return BusinessCollectionTableName
}

type BusinessCollection struct {
	ID          *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string         `gorm:"column:name;not null"`
	BusinessId  *uuid.UUID     `gorm:"column:business_id;not null"`
	Index       int32          `gorm:"column:index;not null"`
	EnabledFlag int32          `gorm:"column:enabled_flag;not null"`
	Item        []Item         `gorm:"foreignKey:BusinessCollectionId"`
	CreateTime  time.Time      `gorm:"column:create_time;not null"`
	UpdateTime  time.Time      `gorm:"column:update_time;not null"`
	DeleteTime  gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessCollection) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *BusinessCollection) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
