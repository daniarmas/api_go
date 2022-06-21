package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UserPermissionTableName = "user_permission"

func (UserPermission) TableName() string {
	return UserPermissionTableName
}

type UserPermission struct {
	ID           *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name         string         `gorm:"column:name;not null"`
	UserId       *uuid.UUID     `gorm:"column:user_id;not null"`
	User         User           `gorm:"foreignKey:UserId"`
	BusinessId   *uuid.UUID     `gorm:"column:business_id"`
	Business     Business       `gorm:"foreignKey:BusinessId"`
	PermissionId *uuid.UUID     `gorm:"column:permission_id;not null"`
	Permission   Permission     `gorm:"foreignKey:PermissionId"`
	CreateTime   time.Time      `gorm:"column:create_time;not null"`
	UpdateTime   time.Time      `gorm:"column:update_time;not null"`
	DeleteTime   gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *UserPermission) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
