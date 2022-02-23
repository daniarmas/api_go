package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UnionOrderAndOrderedItemTableName = "union_order_and_ordered_item"

func (UnionOrderAndOrderedItem) TableName() string {
	return UnionOrderAndOrderedItemTableName
}

type UnionOrderAndOrderedItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	OrderFk       uuid.UUID      `gorm:"column:order_fk;not null"`
	OrderedItemFk uuid.UUID      `gorm:"column:ordered_item_fk;not null"`
	CreateTime    time.Time      `gorm:"column:create_time;not null"`
	UpdateTime    time.Time      `gorm:"column:update_time;not null"`
	DeleteTime    gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UnionOrderAndOrderedItem) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UnionOrderAndOrderedItem) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
