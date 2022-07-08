package datasource

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessScheduleDatasource interface {
	GetBusinessSchedule(tx *gorm.DB, where *entity.BusinessSchedule, fields *[]string) (*entity.BusinessSchedule, error)
	BusinessIsOpen(tx *gorm.DB, where *entity.BusinessSchedule, orderType string) (bool, error)
}

type businessScheduleDatasource struct{}

func (v *businessScheduleDatasource) GetBusinessSchedule(tx *gorm.DB, where *entity.BusinessSchedule, fields *[]string) (*entity.BusinessSchedule, error) {
	var res *entity.BusinessSchedule
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (v *businessScheduleDatasource) BusinessIsOpen(tx *gorm.DB, where *entity.BusinessSchedule, orderType string) (bool, error) {
	var schedule *entity.BusinessSchedule
	timeNow := time.Now().UTC()
	weekday := timeNow.Weekday().String()
	result := tx.Where(where).Take(&schedule)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return false, errors.New("record not found")
		} else {
			return false, result.Error
		}
	}
	if orderType == "OrderTypeHomeDelivery" {
		switch weekday {
		case "Sunday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliverySunday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliverySunday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeSunday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeSunday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeSunday) || timeNow.After(closingTimeSunday) {
				return false, errors.New("business closed")
			}
		case "Monday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliveryMonday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliveryMonday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeMonday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeMonday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeMonday) || timeNow.After(closingTimeMonday) {
				return false, errors.New("business closed")
			}
		case "Tuesday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliveryTuesday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliveryTuesday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeTuesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeTuesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeTuesday) || timeNow.After(closingTimeTuesday) {
				return false, errors.New("business closed")
			}
		case "Wednesday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliveryWednesday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliveryWednesday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeWednesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeWednesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeWednesday) || timeNow.After(closingTimeWednesday) {
				return false, errors.New("business closed")
			}
		case "Thursday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliveryThursday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliveryThursday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeThursday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeThursday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeThursday) || timeNow.After(closingTimeThursday) {
				return false, errors.New("business closed")
			}
		case "Friday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliveryFriday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliveryFriday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeFriday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeFriday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeFriday) || timeNow.After(closingTimeFriday) {
				return false, errors.New("business closed")
			}
		case "Saturday":
			splitOpening := strings.Split(schedule.OpeningTimeDeliverySaturday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeDeliverySaturday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeSaturday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeSaturday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeSaturday) || timeNow.After(closingTimeSaturday) {
				return false, errors.New("business closed")
			}
		}
	} else if orderType == "OrderTypePickUp" {
		switch weekday {
		case "Sunday":
			splitOpening := strings.Split(schedule.OpeningTimeSunday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeSunday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeSunday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeSunday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeSunday) || timeNow.After(closingTimeSunday) {
				return false, errors.New("business closed")
			}
		case "Monday":
			splitOpening := strings.Split(schedule.OpeningTimeMonday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeMonday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeMonday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeMonday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeMonday) || timeNow.After(closingTimeMonday) {
				return false, errors.New("business closed")
			}
		case "Tuesday":
			splitOpening := strings.Split(schedule.OpeningTimeTuesday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeTuesday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeTuesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeTuesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeTuesday) || timeNow.After(closingTimeTuesday) {
				return false, errors.New("business closed")
			}
		case "Wednesday":
			splitOpening := strings.Split(schedule.OpeningTimeWednesday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeWednesday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeWednesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeWednesday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeWednesday) || timeNow.After(closingTimeWednesday) {
				return false, errors.New("business closed")
			}
		case "Thursday":
			splitOpening := strings.Split(schedule.OpeningTimeThursday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeThursday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeThursday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeThursday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeThursday) || timeNow.After(closingTimeThursday) {
				return false, errors.New("business closed")
			}
		case "Friday":
			splitOpening := strings.Split(schedule.OpeningTimeFriday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeFriday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeFriday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeFriday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeFriday) || timeNow.After(closingTimeFriday) {
				return false, errors.New("business closed")
			}
		case "Saturday":
			splitOpening := strings.Split(schedule.OpeningTimeSaturday, ":")
			splitClosing := strings.Split(schedule.ClosingTimeSaturday, ":")
			openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
			if openingHourErr != nil {
				return false, openingHourErr
			}
			openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
			if openingMinutesErr != nil {
				return false, openingMinutesErr
			}
			closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
			if closingHourErr != nil {
				return false, closingHourErr
			}
			closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
			if closingMinutesErr != nil {
				return false, closingMinutesErr
			}
			openingTimeSaturday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
			closingTimeSaturday := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
			if timeNow.Before(openingTimeSaturday) || timeNow.After(closingTimeSaturday) {
				return false, errors.New("business closed")
			}
		}
	}
	return true, nil
}
