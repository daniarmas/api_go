package entity

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
	ID                         *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	FirstOpeningTimeSunday     time.Time      `gorm:"column:first_opening_time_sunday"`
	SecondOpeningTimeSunday    time.Time      `gorm:"column:second_opening_time_sunday"`
	FirstOpeningTimeMonday     time.Time      `gorm:"column:first_opening_time_monday"`
	SecondOpeningTimeMonday    time.Time      `gorm:"column:second_opening_time_monday"`
	FirstOpeningTimeTuesday    time.Time      `gorm:"column:first_opening_time_tuesday"`
	SecondOpeningTimeTuesday   time.Time      `gorm:"column:second_opening_time_tuesday"`
	FirstOpeningTimeWednesday  time.Time      `gorm:"column:first_opening_time_wednesday"`
	SecondOpeningTimeWednesday time.Time      `gorm:"column:second_opening_time_wednesday"`
	FirstOpeningTimeThursday   time.Time      `gorm:"column:first_opening_time_thursday"`
	SecondOpeningTimeThursday  time.Time      `gorm:"column:second_opening_time_thursday"`
	FirstOpeningTimeFriday     time.Time      `gorm:"column:first_opening_time_friday"`
	SecondOpeningTimeFriday    time.Time      `gorm:"column:second_opening_time_friday"`
	FirstOpeningTimeSaturday   time.Time      `gorm:"column:first_opening_time_saturday"`
	SecondOpeningTimeSaturday  time.Time      `gorm:"column:second_opening_time_saturday"`
	FirstClosingTimeSunday     time.Time      `gorm:"column:first_closing_time_sunday"`
	SecondClosingTimeSunday    time.Time      `gorm:"column:second_closing_time_sunday"`
	FirstClosingTimeMonday     time.Time      `gorm:"column:first_closing_time_monday"`
	SecondClosingTimeMonday    time.Time      `gorm:"column:second_closing_time_monday"`
	FirstClosingTimeTuesday    time.Time      `gorm:"column:first_closing_time_tuesday"`
	SecondClosingTimeTuesday   time.Time      `gorm:"column:second_closing_time_tuesday"`
	FirstClosingTimeWednesday  time.Time      `gorm:"column:first_closing_time_wednesday"`
	SecondClosingTimeWednesday time.Time      `gorm:"column:second_closing_time_wednesday"`
	FirstClosingTimeThursday   time.Time      `gorm:"column:first_closing_time_thursday"`
	SecondClosingTimeThursday  time.Time      `gorm:"column:second_closing_time_thursday"`
	FirstClosingTimeFriday     time.Time      `gorm:"column:first_closing_time_friday"`
	SecondClosingTimeFriday    time.Time      `gorm:"column:second_closing_time_friday"`
	FirstClosingTimeSaturday   time.Time      `gorm:"column:first_closing_time_saturday"`
	SecondClosingTimeSaturday  time.Time      `gorm:"column:second_closing_time_saturday"`
	BusinessId                 *uuid.UUID     `gorm:"column:business_id;not null"`
	CreateTime                 time.Time      `gorm:"column:create_time;not null"`
	UpdateTime                 time.Time      `gorm:"column:update_time;not null"`
	DeleteTime                 gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *BusinessSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	i.UpdateTime = time.Now().UTC()
	return
}
