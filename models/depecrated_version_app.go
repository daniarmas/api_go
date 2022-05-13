package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const DeprecatedVersionAppTableName = "deprecated_version_app"

func (DeprecatedVersionApp) TableName() string {
	return DeprecatedVersionAppTableName
}

type DeprecatedVersionApp struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Version     string         `gorm:"column:version;not null"`
	Description string         `gorm:"column:description;not null"`
	CreateTime  time.Time      `gorm:"column:create_time;not null"`
	UpdateTime  time.Time      `gorm:"column:update_time;not null"`
	DeleteTime  gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *DeprecatedVersionApp) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *DeprecatedVersionApp) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
