package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

const OrderTableName = "order"

func (Order) TableName() string {
	return OrderTableName
}

type Order struct {
	ID                   *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	ShortId              string         `gorm:"column:short_id;not null"`
	Status               string         `gorm:"column:status"`
	ItemsQuantity        int32          `gorm:"column:items_quantity;not null"`
	OrderType            string         `gorm:"column:order_type;not null"`
	DeliveryPriceCup     string         `gorm:"column:delivery_price_cup;not null"`
	PriceCup             string         `gorm:"column:price_cup;not null"`
	Number               string         `gorm:"column:number;not null"`
	Address              string         `gorm:"column:address;not null"`
	Instructions         string         `gorm:"column:instructions"`
	Phone                string         `gorm:"column:phone"`
	PaymentMethodType    string         `gorm:"column:payment_method_type"`
	CancelReasons        string         `gorm:"column:cancel_reasons"`
	Coordinates          ewkb.Point     `gorm:"column:coordinates;not null"`
	AuthorizationTokenId *uuid.UUID     `gorm:"column:authorization_token_id;not null"`
	UserId               *uuid.UUID     `gorm:"column:user_id;not null"`
	User                 User           `gorm:"foreignKey:UserId"`
	BusinessId           *uuid.UUID     `gorm:"column:business_id;not null"`
	BusinessThumbnail    string         `gorm:"column:business_thumbnail;not null"`
	BusinessName         string         `gorm:"column:business_name;not null"`
	Business             Business       `gorm:"foreignKey:BusinessId"`
	StartOrderTime       time.Time      `gorm:"column:start_order_time;not null"`
	EndOrderTime         time.Time      `gorm:"column:end_order_time;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}
type OrderBusiness struct {
	ID                   *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	ShortId              string         `gorm:"column:short_id;not null"`
	BusinessName         string         `gorm:"column:business_name"`
	Status               string         `gorm:"column:status"`
	ItemsQuantity        int32          `gorm:"column:items_quantity;not null"`
	OrderType            string         `gorm:"column:order_type;not null"`
	DeliveryPriceCup     string         `gorm:"column:delivery_price_cup;not null"`
	PriceCup             string         `gorm:"column:price_cup;not null"`
	Number               string         `gorm:"column:number;not null"`
	Address              string         `gorm:"column:address;not null"`
	Instructions         string         `gorm:"column:instructions"`
	CancelReasons        string         `gorm:"column:cancel_reasons"`
	Coordinates          ewkb.Point     `gorm:"column:coordinates;not null"`
	AuthorizationTokenId *uuid.UUID     `gorm:"column:authorization_token_id;not null"`
	UserId               *uuid.UUID     `gorm:"column:user_id;not null"`
	User                 User           `gorm:"foreignKey:UserId"`
	BusinessId           *uuid.UUID     `gorm:"column:business_id;not null"`
	BusinessThumbnail    string         `gorm:"column:business_thumbnail;not null"`
	Business             Business       `gorm:"foreignKey:BusinessId"`
	StartOrderTime       time.Time      `gorm:"column:start_order_time;not null"`
	EndOrderTime         time.Time      `gorm:"column:end_order_time;not null"`
	CreateTime           time.Time      `gorm:"column:create_time;not null"`
	UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *Order) BeforeCreate(tx *gorm.DB) (err error) {
	var shortId string
	if shortId, err = shortid.Generate(); err != nil {
		return err
	}
	i.ShortId = shortId
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}

func (u *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateTime = time.Now().UTC()
	return
}
