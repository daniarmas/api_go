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
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Tag            string         `gorm:"column:tag;not null"`
	UserFk         uuid.UUID      `gorm:"column:user_fk;not null"`
	Coordinates    ewkb.Point     `gorm:"column:coordinates;not null"`
	ResidenceType  string         `gorm:"column:residence_type;not null"`
	HouseNumber    string         `gorm:"column:house_number"`
	BuildingNumber string         `gorm:"column:building_number"`
	Description    string         `gorm:"column:description"`
	ProvinceFk     uuid.UUID      `gorm:"column:province_fk;not null"`
	MunicipalityFk uuid.UUID      `gorm:"column:municipality_fk;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UserAddress) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	return
}

func (u *UserAddress) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now()
	return
}
