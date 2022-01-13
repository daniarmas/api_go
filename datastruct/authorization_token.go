package datastruct

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

func (AuthorizationToken) TableName() string {
	return AuthorizationTokenTableName
}

const AuthorizationTokenTableName = "authorization_token"

type AuthorizationToken struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	RefreshTokenFk uuid.UUID      `gorm:"type:uuid;column:refresh_token_fk"`
	UserFk         uuid.UUID      `gorm:"column:user_fk"`
	DeviceFk       uuid.UUID      `gorm:"column:device_fk"`
	App            string         `gorm:"column:app"`
	AppVersion     string         `gorm:"column:app_version"`
	CreateTime     time.Time      `gorm:"column:create_time"`
	UpdateTime     time.Time      `gorm:"column:update_time"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *AuthorizationToken) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
