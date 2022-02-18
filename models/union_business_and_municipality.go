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
	BusinessFk     uuid.UUID      `gorm:"column:business_fk;not null"`
	MunicipalityFk uuid.UUID      `gorm:"column:municipality_fk;not null"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *UnionBusinessAndMunicipality) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	return
}

func (u *UnionBusinessAndMunicipality) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now()
	return
}