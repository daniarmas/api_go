package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const OrderTableName = "order"

func (Order) TableName() string {
	return OrderTableName
}

type Order struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Status               string         `gorm:"column:status"`
	Quantity             int32          `gorm:"column:quantity"`
	DeliveryType         string         `gorm:"column:delivery_type"`
	ResidenceType        string         `gorm:"column:residence_type"`
	Price                float64        `gorm:"column:price"`
	BuildingNumber       string         `gorm:"column:building_number"`
	HouseNumber          string         `gorm:"column:house_number"`
	BusinessFk           uuid.UUID      `gorm:"column:business_fk;not null"`
	Business             Business       `gorm:"foreignKey:BusinessFk"`
	Coordinates          ewkb.Point     `gorm:"column:coordinates"`
	UserFk               uuid.UUID      `gorm:"column:user_fk;not null"`
	AuthorizationTokenFk uuid.UUID      `gorm:"column:authorization_token_fk;not null"`
	DeliveryDate         time.Time      `gorm:"column:delivery_date;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}
type OrderBusiness struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	BusinessName         string         `gorm:"column:business_name"`
	Quantity             int32          `gorm:"column:quantity"`
	Status               string         `gorm:"column:status"`
	DeliveryType         string         `gorm:"column:delivery_type"`
	ResidenceType        string         `gorm:"column:residence_type"`
	Price                float64        `gorm:"column:price"`
	BuildingNumber       string         `gorm:"column:building_number"`
	HouseNumber          string         `gorm:"column:house_number"`
	BusinessFk           uuid.UUID      `gorm:"column:business_fk;not null"`
	Business             Business       `gorm:"foreignKey:BusinessFk"`
	Coordinates          ewkb.Point     `gorm:"column:coordinates"`
	UserFk               uuid.UUID      `gorm:"column:user_fk;not null"`
	AuthorizationTokenFk uuid.UUID      `gorm:"column:authorization_token_fk;not null"`
	DeliveryDate         time.Time      `gorm:"column:delivery_date;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *Order) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
