package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UnionBusinessAndMunicipalityTableName = "union_business_and_municipality"

func (UnionBusinessAndMunicipality) TableName() string {
	return UnionBusinessAndMunicipalityTableName
}

type UnionBusinessAndMunicipality struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessId     uuid.UUID      `gorm:"column:business_id;not null"`
	MunicipalityId uuid.UUID      `gorm:"column:municipality_id;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}
type UnionBusinessAndMunicipalityWithMunicipality struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessId       uuid.UUID      `gorm:"column:business_id;not null"`
	MunicipalityId   uuid.UUID      `gorm:"column:municipality_id;not null"`
	MunicipalityName string         `gorm:"column:municipality_name;not null"`
	CreateTime       time.Time      `gorm:"column:create_time;not null"`
	UpdateTime       time.Time      `gorm:"column:update_time;not null"`
	DeleteTime       gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UnionBusinessAndMunicipality) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC()
	u.UpdateTime = time.Now().UTC()
	return
}

func (u *UnionBusinessAndMunicipality) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
