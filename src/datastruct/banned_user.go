package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BannedUserTableName = "item"

type BannedUser struct {
	ID                            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Description                   string         `gorm:"column:description"`
	UserFk                        float64        `gorm:"column:userFk"`
	Email                         int64          `gorm:"column:email"`
	ModeratorAuthorizationTokenFk string         `gorm:"column:moderatorAuthorizationTokenFk"`
	CreateTime                    time.Time      `gorm:"column:createTime"`
	UpdateTime                    time.Time      `gorm:"column:updateTime"`
	DeleteTime                    gorm.DeletedAt `gorm:"index;column:deleteTime"`
}
