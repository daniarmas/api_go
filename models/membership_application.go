package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const MembershipApplicationTableName = "membership_application"

func (MembershipApplication) TableName() string {
	return MembershipApplicationTableName
}

type MembershipApplication struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessName   string         `gorm:"column:business_name"`
	Description    string         `gorm:"column:description"`
	Coordinates    ewkb.Point     `gorm:"column:coordinates"`
	Status         string         `gorm:"column:status"`
	ProvinceId     *uuid.UUID     `gorm:"column:province_id"`
	MunicipalityId *uuid.UUID     `gorm:"column:municipality_id"`
	CreateTime     time.Time      `gorm:"column:create_time"`
	UpdateTime     time.Time      `gorm:"column:update_time"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *MembershipApplication) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *MembershipApplication) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
