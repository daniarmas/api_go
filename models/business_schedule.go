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
	ID                           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	OpeningTimeSunday            time.Time      `gorm:"column:opening_time_sunday"`
	OpeningTimeMonday            time.Time      `gorm:"column:opening_time_monday"`
	OpeningTimeTuesday           time.Time      `gorm:"column:opening_time_tuesday"`
	OpeningTimeWednesday         time.Time      `gorm:"column:opening_time_wednesday"`
	OpeningTimeThursday          time.Time      `gorm:"column:opening_time_thursday"`
	OpeningTimeFriday            time.Time      `gorm:"column:opening_time_friday"`
	OpeningTimeSaturday          time.Time      `gorm:"column:opening_time_saturday"`
	ClosingTimeSunday            time.Time      `gorm:"column:closing_time_sunday"`
	ClosingTimeMonday            time.Time      `gorm:"column:closing_time_monday"`
	ClosingTimeTuesday           time.Time      `gorm:"column:closing_time_tuesday"`
	ClosingTimeWednesday         time.Time      `gorm:"column:closing_time_wednesday"`
	ClosingTimeThursday          time.Time      `gorm:"column:closing_time_thursday"`
	ClosingTimeFriday            time.Time      `gorm:"column:closing_time_friday"`
	ClosingTimeSaturday          time.Time      `gorm:"column:closing_time_saturday"`
	OpeningTimeDeliverySunday    time.Time      `gorm:"column:opening_time_delivery_sunday"`
	OpeningTimeDeliveryMonday    time.Time      `gorm:"column:opening_time_delivery_monday"`
	OpeningTimeDeliveryTuesday   time.Time      `gorm:"column:opening_time_delivery_tuesday"`
	OpeningTimeDeliveryWednesday time.Time      `gorm:"column:opening_time_delivery_wednesday"`
	OpeningTimeDeliveryThursday  time.Time      `gorm:"column:opening_time_delivery_thursday"`
	OpeningTimeDeliveryFriday    time.Time      `gorm:"column:opening_time_delivery_friday"`
	OpeningTimeDeliverySaturday  time.Time      `gorm:"column:opening_time_delivery_saturday"`
	ClosingTimeDeliverySunday    time.Time      `gorm:"column:opening_time_delivery_sunday"`
	ClosingTimeDeliveryMonday    time.Time      `gorm:"column:opening_time_delivery_monday"`
	ClosingTimeDeliveryTuesday   time.Time      `gorm:"column:opening_time_delivery_tuesday"`
	ClosingTimeDeliveryWednesday time.Time      `gorm:"column:opening_time_delivery_wednesday"`
	ClosingTimeDeliveryThursday  time.Time      `gorm:"column:opening_time_delivery_thursday"`
	ClosingTimeDeliveryFriday    time.Time      `gorm:"column:opening_time_delivery_friday"`
	ClosingTimeDeliverySaturday  time.Time      `gorm:"column:opening_time_delivery_saturday"`
	BusinessFk                   uuid.UUID      `gorm:"column:business_fk;not null"`
	CreateTime                   time.Time      `gorm:"column:create_time;not null"`
	UpdateTime                   time.Time      `gorm:"column:update_time;not null"`
	DeleteTime                   gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}
