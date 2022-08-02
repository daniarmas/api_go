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
	FirstOpeningTimeSunday     string         `gorm:"column:first_opening_time_sunday"`
	SecondOpeningTimeSunday    string         `gorm:"column:second_opening_time_sunday"`
	FirstOpeningTimeMonday     string         `gorm:"column:first_opening_time_monday"`
	SecondOpeningTimeMonday    string         `gorm:"column:second_opening_time_monday"`
	FirstOpeningTimeTuesday    string         `gorm:"column:first_opening_time_tuesday"`
	SecondOpeningTimeTuesday   string         `gorm:"column:second_opening_time_tuesday"`
	FirstOpeningTimeWednesday  string         `gorm:"column:first_opening_time_wednesday"`
	SecondOpeningTimeWednesday string         `gorm:"column:second_opening_time_wednesday"`
	FirstOpeningTimeThursday   string         `gorm:"column:first_opening_time_thursday"`
	SecondOpeningTimeThursday  string         `gorm:"column:second_opening_time_thursday"`
	FirstOpeningTimeFriday     string         `gorm:"column:first_opening_time_friday"`
	SecondOpeningTimeFriday    string         `gorm:"column:second_opening_time_friday"`
	FirstOpeningTimeSaturday   string         `gorm:"column:first_opening_time_saturday"`
	SecondOpeningTimeSaturday  string         `gorm:"column:second_opening_time_saturday"`
	FirstClosingTimeSunday     string         `gorm:"column:first_closing_time_sunday"`
	SecondClosingTimeSunday    string         `gorm:"column:second_closing_time_sunday"`
	FirstClosingTimeMonday     string         `gorm:"column:first_closing_time_monday"`
	SecondClosingTimeMonday    string         `gorm:"column:second_closing_time_monday"`
	FirstClosingTimeTuesday    string         `gorm:"column:first_closing_time_tuesday"`
	SecondClosingTimeTuesday   string         `gorm:"column:second_closing_time_tuesday"`
	FirstClosingTimeWednesday  string         `gorm:"column:first_closing_time_wednesday"`
	SecondClosingTimeWednesday string         `gorm:"column:second_closing_time_wednesday"`
	FirstClosingTimeThursday   string         `gorm:"column:first_closing_time_thursday"`
	SecondClosingTimeThursday  string         `gorm:"column:second_closing_time_thursday"`
	FirstClosingTimeFriday     string         `gorm:"column:first_closing_time_friday"`
	SecondClosingTimeFriday    string         `gorm:"column:second_closing_time_friday"`
	FirstClosingTimeSaturday   string         `gorm:"column:first_closing_time_saturday"`
	SecondClosingTimeSaturday  string         `gorm:"column:second_closing_time_saturday"`
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
