package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const PartnerApplicationTableName = "partner_application"

func (PartnerApplication) TableName() string {
	return PartnerApplicationTableName
}

type PartnerApplication struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessName   string         `gorm:"column:business_name"`
	Description    string         `gorm:"column:description"`
	Coordinates    ewkb.Point     `gorm:"column:coordinates"`
	Status         string         `gorm:"column:status"`
	ProvinceId     *uuid.UUID     `gorm:"column:province_id"`
	MunicipalityId *uuid.UUID     `gorm:"column:municipality_id"`
	UserId         *uuid.UUID     `gorm:"column:user_id"`
	CreateTime     time.Time      `gorm:"column:create_time"`
	UpdateTime     time.Time      `gorm:"column:update_time"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *PartnerApplication) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *PartnerApplication) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
