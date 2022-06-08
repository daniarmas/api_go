package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TermTableName = "term"

func (Term) TableName() string {
	return TermTableName
}

type Term struct {
	ID              *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name            string         `gorm:"column:name;not null"`
	Description     string         `gorm:"column:description;not null"`
	BanTimeInYears  int            `gorm:"column:ban_time_in_years;not null"`
	BanTimeInMonths int            `gorm:"column:ban_time_in_months;not null"`
	BanTimeInDays   int            `gorm:"column:ban_time_in_days;not null"`
	BanTimeInHours  int            `gorm:"column:ban_time_in_hours;not null"`
	CreateTime      time.Time      `gorm:"column:create_time;not null"`
	UpdateTime      time.Time      `gorm:"column:update_time;not null"`
	DeleteTime      gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *Term) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *Term) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
