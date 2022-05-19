package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const UserAddressTableName = "user_address"

func (UserAddress) TableName() string {
	return UserAddressTableName
}

type UserAddress struct {
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Tag            string         `gorm:"column:tag;not null"`
	UserId         *uuid.UUID     `gorm:"column:user_id;not null"`
	User           User           `gorm:"foreignKey:UserId"`
	Coordinates    ewkb.Point     `gorm:"column:coordinates;not null"`
	ResidenceType  string         `gorm:"column:residence_type;not null"`
	Address        string         `gorm:"column:address;not null"`
	Number         string         `gorm:"column:number;not null"`
	Instructions   string         `gorm:"column:instructions"`
	ProvinceId     *uuid.UUID     `gorm:"column:province_id;not null"`
	MunicipalityId *uuid.UUID     `gorm:"column:municipality_id;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UserAddress) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UserAddress) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
