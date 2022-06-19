package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const PermissionTableName = "permission"

func (Permission) TableName() string {
	return PermissionTableName
}

type Permission struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string         `gorm:"column:name;not null"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
