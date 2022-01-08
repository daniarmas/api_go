package datastruct

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
	ID                            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string         `gorm:"column:description"`
	UserFk                        uuid.UUID      `gorm:"column:user_fk"`
	Email                         string         `gorm:"column:email"`
	ModeratorAuthorizationTokenFk uuid.UUID      `gorm:"column:moderator_authorization_token_fk"`
	CreateTime                    time.Time      `gorm:"column:create_time"`
	UpdateTime                    time.Time      `gorm:"column:update_time"`
	DeleteTime                    gorm.DeletedAt `gorm:"index;column:delete_time"`
}
