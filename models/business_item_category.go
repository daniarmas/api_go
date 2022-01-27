package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessItemCategoryTableName = "business_item_category"

func (BusinessItemCategory) TableName() string {
	return BusinessItemCategoryTableName
}

type BusinessItemCategory struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string         `gorm:"column:name;not null"`
	BusinessFk uuid.UUID      `gorm:"column:business_fk;not null"`
	Index      int32          `gorm:"column:index;not null"`
	Item       []Item         `gorm:"foreignKey:BusinessItemCategoryFk"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessItemCategory) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}

func (u *BusinessItemCategory) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now()
	return
}
