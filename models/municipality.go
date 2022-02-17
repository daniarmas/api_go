package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const MunicipalityTableName = "municipality"

func (Municipality) TableName() string {
	return MunicipalityTableName
}

type Municipality struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string         `gorm:"column:name;not null"`
	Zoom        float32        `gorm:"column:zoom;not null"`
	Coordinates ewkb.Point     `gorm:"column:coordinates"`
	Polygon     ewkb.Polygon   `gorm:"column:polygon"`
	ProvinceFk  uuid.UUID      `gorm:"column:province_fk;not null"`
	CreateTime  time.Time      `gorm:"column:create_time;not null"`
	UpdateTime  time.Time      `gorm:"column:update_time;not null"`
	DeleteTime  gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *Municipality) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now()
	r.UpdateTime = time.Now()
	return
}
