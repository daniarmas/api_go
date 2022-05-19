package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const BusinessScheduleTableName = "business_schedule"

func (BusinessSchedule) TableName() string {
	return BusinessScheduleTableName
}

type BusinessSchedule struct {
	ID                           *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	OpeningTimeSunday            string         `gorm:"column:opening_time_sunday"`
	OpeningTimeMonday            string         `gorm:"column:opening_time_monday"`
	OpeningTimeTuesday           string         `gorm:"column:opening_time_tuesday"`
	OpeningTimeWednesday         string         `gorm:"column:opening_time_wednesday"`
	OpeningTimeThursday          string         `gorm:"column:opening_time_thursday"`
	OpeningTimeFriday            string         `gorm:"column:opening_time_friday"`
	OpeningTimeSaturday          string         `gorm:"column:opening_time_saturday"`
	ClosingTimeSunday            string         `gorm:"column:closing_time_sunday"`
	ClosingTimeMonday            string         `gorm:"column:closing_time_monday"`
	ClosingTimeTuesday           string         `gorm:"column:closing_time_tuesday"`
	ClosingTimeWednesday         string         `gorm:"column:closing_time_wednesday"`
	ClosingTimeThursday          string         `gorm:"column:closing_time_thursday"`
	ClosingTimeFriday            string         `gorm:"column:closing_time_friday"`
	ClosingTimeSaturday          string         `gorm:"column:closing_time_saturday"`
	OpeningTimeDeliverySunday    string         `gorm:"column:opening_time_delivery_sunday"`
	OpeningTimeDeliveryMonday    string         `gorm:"column:opening_time_delivery_monday"`
	OpeningTimeDeliveryTuesday   string         `gorm:"column:opening_time_delivery_tuesday"`
	OpeningTimeDeliveryWednesday string         `gorm:"column:opening_time_delivery_wednesday"`
	OpeningTimeDeliveryThursday  string         `gorm:"column:opening_time_delivery_thursday"`
	OpeningTimeDeliveryFriday    string         `gorm:"column:opening_time_delivery_friday"`
	OpeningTimeDeliverySaturday  string         `gorm:"column:opening_time_delivery_saturday"`
	ClosingTimeDeliverySunday    string         `gorm:"column:closing_time_delivery_sunday"`
	ClosingTimeDeliveryMonday    string         `gorm:"column:closing_time_delivery_monday"`
	ClosingTimeDeliveryTuesday   string         `gorm:"column:closing_time_delivery_tuesday"`
	ClosingTimeDeliveryWednesday string         `gorm:"column:closing_time_delivery_wednesday"`
	ClosingTimeDeliveryThursday  string         `gorm:"column:closing_time_delivery_thursday"`
	ClosingTimeDeliveryFriday    string         `gorm:"column:closing_time_delivery_friday"`
	ClosingTimeDeliverySaturday  string         `gorm:"column:closing_time_delivery_saturday"`
	BusinessId                   *uuid.UUID     `gorm:"column:business_id;not null"`
	CreateTime                   time.Time      `gorm:"column:create_time;not null"`
	UpdateTime                   time.Time      `gorm:"column:update_time;not null"`
	DeleteTime                   gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}
