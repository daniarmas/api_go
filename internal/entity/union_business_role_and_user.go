package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UnionBusinessRoleAndUserTableName = "union_business_role_and_user"

func (UnionBusinessRoleAndUser) TableName() string {
	return UnionBusinessRoleAndUserTableName
}

type UnionBusinessRoleAndUser struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessRoleId *uuid.UUID     `gorm:"column:business_role_id;not null"`
	UserId         *uuid.UUID     `gorm:"column:user_id;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UnionBusinessRoleAndUser) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UnionBusinessRoleAndUser) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
