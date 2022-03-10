package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BannedUserTableName = "banned_user"

func (BannedUser) TableName() string {
	return BannedUserTableName
}

type BannedUser struct {
	ID                            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string             `gorm:"column:description;not null"`
	UserFk                        uuid.UUID          `gorm:"column:user_fk;not null"`
	User                          User               `gorm:"foreignKey:UserFk"`
	Email                         string             `gorm:"column:email;not null"`
	ModeratorAuthorizationTokenFk uuid.UUID          `gorm:"column:moderator_authorization_token_fk;not null"`
	AuthorizationToken            AuthorizationToken `gorm:"foreignKey:ModeratorAuthorizationTokenFk"`
	BanExpirationTime             time.Time          `gorm:"column:ban_expiration_time;not null"`
	CreateTime                    time.Time          `gorm:"column:create_time;not null"`
	UpdateTime                    time.Time          `gorm:"column:update_time;not null"`
	DeleteTime                    gorm.DeletedAt     `gorm:"index;column:delete_time"`
}

func (i *BannedUser) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}
