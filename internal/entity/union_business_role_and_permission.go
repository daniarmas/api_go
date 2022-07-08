package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UnionBusinessRoleAndPermissionTableName = "union_business_role_and_permission"

func (UnionBusinessRoleAndPermission) TableName() string {
	return UnionBusinessRoleAndPermissionTableName
}

type UnionBusinessRoleAndPermission struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessRoleId *uuid.UUID     `gorm:"column:business_role_id;not null"`
	PermissionId   *uuid.UUID     `gorm:"column:permission_id;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}
type UnionBusinessRoleAndPermissionWithPermission struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name           string         `gorm:"column:name;not null"`
	BusinessRoleId *uuid.UUID     `gorm:"column:business_role_id;not null"`
	PermissionId   *uuid.UUID     `gorm:"column:permission_id;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UnionBusinessRoleAndPermission) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UnionBusinessRoleAndPermission) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
