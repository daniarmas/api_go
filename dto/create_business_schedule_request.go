package dto

import "github.com/google/uuid"

type CreateBusinessScheduleRequest struct {
	BusinessFk                   uuid.UUID
	OpeningTimeSunday            string
	OpeningTimeMonday            string
	OpeningTimeTuesday           string
	OpeningTimeWednesday         string
	OpeningTimeThursday          string
	OpeningTimeFriday            string
	OpeningTimeSaturday          string
	ClosingTimeSunday            string
	ClosingTimeMonday            string
	ClosingTimeTuesday           string
	ClosingTimeWednesday         string
	ClosingTimeThursday          string
	ClosingTimeFriday            string
	ClosingTimeSaturday          string
	OpeningTimeDeliverySunday    string
	OpeningTimeDeliveryMonday    string
	OpeningTimeDeliveryTuesday   string
	OpeningTimeDeliveryWednesday string
	OpeningTimeDeliveryThursday  string
	OpeningTimeDeliveryFriday    string
	OpeningTimeDeliverySaturday  string
	ClosingTimeDeliverySunday    string
	ClosingTimeDeliveryMonday    string
	ClosingTimeDeliveryTuesday   string
	ClosingTimeDeliveryWednesday string
	ClosingTimeDeliveryThursday  string
	ClosingTimeDeliveryFriday    string
	ClosingTimeDeliverySaturday  string
}
