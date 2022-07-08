package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessRoleTableName = "business_role"

func (BusinessRole) TableName() string {
	return BusinessRoleTableName
}

type BusinessRole struct {
	ID         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string         `gorm:"column:name;not null"`
	BusinessId *uuid.UUID     `gorm:"column:business_id;not null"`
	Business   Business       `gorm:"foreignKey:BusinessId"`
	Permission []Permission   `gorm:"many2many:union_business_role_and_permission;"`
	CreateTime time.Time      `gorm:"column:create_time;not null"`
	UpdateTime time.Time      `gorm:"column:update_time;not null"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessRole) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *BusinessRole) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
