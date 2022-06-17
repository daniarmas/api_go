package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessUserPermissionTableName = "business_user_permission"

func (BusinessUserPermission) TableName() string {
	return BusinessUserPermissionTableName
}

type BusinessUserPermission struct {
	ID                   *uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                 string             `gorm:"column:name;not null"`
	UserId               *uuid.UUID         `gorm:"column:user_id;not null"`
	User                 User               `gorm:"foreignKey:UserId"`
	BusinessId           *uuid.UUID         `gorm:"column:business_id"`
	Business             Business           `gorm:"foreignKey:BusinessId"`
	BusinessPermissionId *uuid.UUID         `gorm:"column:business_permission_id;not null"`
	BusinessPermission   BusinessPermission `gorm:"foreignKey:BusinessPermissionId"`
	CreateTime           time.Time          `gorm:"column:create_time;not null"`
	UpdateTime           time.Time          `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt     `gorm:"index;column:delete_time"`
}

func (r *BusinessUserPermission) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
