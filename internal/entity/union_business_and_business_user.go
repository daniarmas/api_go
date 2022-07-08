package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UnionBusinessAndBusinessUserTableName = "union_business_and_business_user"

func (UnionBusinessAndBusinessUser) TableName() string {
	return UnionBusinessAndBusinessUserTableName
}

type UnionBusinessAndBusinessUser struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessId     *uuid.UUID     `gorm:"column:business_id;not null"`
	BusinessUserId *uuid.UUID     `gorm:"column:business_user_id;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UnionBusinessAndBusinessUser) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UnionBusinessAndBusinessUser) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
