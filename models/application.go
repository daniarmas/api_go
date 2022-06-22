package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ApplicationTableName = "application"

func (Application) TableName() string {
	return ApplicationTableName
}

type Application struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name           string         `gorm:"column:name;not null"`
	Version        string         `gorm:"column:version;not null"`
	Description    string         `gorm:"column:description"`
	ExpirationTime time.Time      `gorm:"column:expiration_time;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *Application) BeforeCreate(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC()
	i.CreateTime = timeNow
	i.UpdateTime = timeNow
	return
}

func (u *Application) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
