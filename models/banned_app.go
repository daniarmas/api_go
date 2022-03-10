package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BannedAppTableName = "banned_app"

func (BannedApp) TableName() string {
	return BannedAppTableName
}

type BannedApp struct {
	ID                            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string         `gorm:"column:description;not null"`
	Version                       string         `gorm:"column:version;not null"`
	ModeratorAuthorizationTokenFk uuid.UUID      `gorm:"column:moderator_authorization_token_fk;not null"`
	CreateTime                    time.Time      `gorm:"column:create_time;not null"`
	UpdateTime                    time.Time      `gorm:"column:update_time;not null"`
	DeleteTime                    gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BannedApp) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}
