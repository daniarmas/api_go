package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TermsViolatedBannedUserTableName = "terms_violated_banned_user"

func (TermsViolatedBannedUser) TableName() string {
	return TermsViolatedBannedUserTableName
}

type TermsViolatedBannedUser struct {
	ID           *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	BannedUserId *uuid.UUID     `gorm:"column:banned_user_id;not null"`
	BannedUser   BannedUser     `gorm:"foreignKey:BannedUserId"`
	TermId       uuid.UUID      `gorm:"column:term_id;not null"`
	Term         Term           `gorm:"foreignKey:TermId"`
	CreateTime   time.Time      `gorm:"column:create_time;not null"`
	UpdateTime   time.Time      `gorm:"column:update_time;not null"`
	DeleteTime   gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *TermsViolatedBannedUser) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *TermsViolatedBannedUser) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
