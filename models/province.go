package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ProvinceTableName = "province"

func (Province) TableName() string {
	return ProvinceTableName
}

type Province struct {
	ID                       *uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                     string         `gorm:"column:name;not null"`
	Codename                 string         `gorm:"column:codename;not null"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *Province) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now().UTC()
	r.UpdateTime = time.Now().UTC()
	return
}
