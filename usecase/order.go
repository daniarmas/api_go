package usecase

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	ListOrder(request *dto.ListOrderRequest) (*dto.ListOrderResponse, error)
	CreateOrder(request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	UpdateOrder(request *dto.UpdateOrderRequest) (*dto.UpdateOrderResponse, error)
	ListOrderedItemWithItem(request *dto.ListOrderedItemRequest) (*dto.ListOrderedItemResponse, error)
}

type orderService struct {
	dao repository.DAO
}

func NewOrderService(dao repository.DAO) OrderService {
	return &orderService{dao: dao}
}

func (i *orderService) ListOrderedItemWithItem(request *dto.ListOrderedItemRequest) (*dto.ListOrderedItemResponse, error) {
	var response dto.ListOrderedItemResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderFk: request.OrderFk})
		if unionOrderAndOrderedItemErr != nil {
			return unionOrderAndOrderedItemErr
		}
		orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
		for _, item := range *unionOrderAndOrderedItemRes {
			orderedItemFks = append(orderedItemFks, item.OrderedItemFk)
		}
		orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks)
		if orderedItemsErr != nil {
			return orderedItemsErr
		}
		response.OrderedItems = orderedItemsRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) UpdateOrder(request *dto.UpdateOrderRequest) (*dto.UpdateOrderResponse, error) {
	var response dto.UpdateOrderResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if request.Status != "OrderStatusTypeCanceled" {
			return errors.New("invalid status value")
		}
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		if request.Status == "OrderStatusTypeCanceled" {
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderFk: request.Id})
			if unionOrderAndOrderedItemErr != nil {
				return unionOrderAndOrderedItemErr
			}
			orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
			for _, item := range *unionOrderAndOrderedItemRes {
				orderedItemFks = append(orderedItemFks, item.OrderedItemFk)
			}
			orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks)
			if orderedItemsErr != nil {
				return orderedItemsErr
			}
			itemFks := make([]uuid.UUID, 0, len(*orderedItemsRes))
			for _, item := range *orderedItemsRes {
				itemFks = append(itemFks, item.ItemFk)
			}
			itemsRes, itemsErr := i.dao.NewItemQuery().ListItemInIds(tx, itemFks)
			if itemsErr != nil {
				return itemsErr
			}
			for _, item := range *orderedItemsRes {
				var index = -1
				for i, n := range *itemsRes {
					if n.ID == item.ItemFk {
						index = i
					}
				}
				(*itemsRes)[index].Availability += int64(item.Quantity)
			}
			for _, item := range *itemsRes {
				_, updateItemsErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &item)
				if updateItemsErr != nil {
					return updateItemsErr
				}
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserFk: authorizationTokenRes.UserFk}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BuildingNumber: updateOrderRes.BuildingNumber, HouseNumber: updateOrderRes.HouseNumber, BusinessFk: updateOrderRes.BusinessFk, UserFk: updateOrderRes.UserFk, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenFk: updateOrderRes.AuthorizationTokenFk, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) CreateOrder(request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	var response dto.CreateOrderResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		createTime := time.Now().UTC()
		weekday := request.OrderDate.Weekday().String()
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listCartItemRes, listCartItemErr := i.dao.NewCartItemRepository().ListCartItemInIds(tx, *request.CartItems)
		if listCartItemErr != nil {
			return listCartItemErr
		}
		var price float64 = 0
		var quantity int32 = 0
		orderedItems := make([]models.OrderedItem, 0, len(*listCartItemRes))
		for _, item := range *listCartItemRes {
			price += item.Price
			quantity += item.Quantity
			orderedItems = append(orderedItems, models.OrderedItem{Quantity: item.Quantity, Price: item.Price, CartItemFk: item.ID, UserFk: item.UserFk, ItemFk: item.ItemFk})
		}
		businessScheduleRes, businessScheduleErr := i.dao.NewBusinessScheduleRepository().GetBusinessSchedule(tx, &models.BusinessSchedule{BusinessFk: (*listCartItemRes)[0].BusinessFk})
		if businessScheduleErr != nil {
			return businessScheduleErr
		}
		businessRes, businessErr := i.dao.NewBusinessQuery().GetBusinessWithLocation(tx, &models.Business{ID: (*listCartItemRes)[0].BusinessFk, Coordinates: request.Coordinates})
		if businessErr != nil {
			return businessErr
		}
		previousTime := createTime
		previousTime = previousTime.AddDate(0, int(businessRes.TimeMarginOrderMonth), int(businessRes.TimeMarginOrderDay))
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderHour) * time.Hour)
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderMinute) * time.Minute)
		if request.OrderDate.Before(previousTime) {
			return errors.New("invalid schedule")
		}
		if request.OrderType == "OrderTypeHomeDelivery" {
			switch weekday {
			case "Sunday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliverySunday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliverySunday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeSunday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSunday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeSunday) || request.OrderDate.After(closingTimeSunday) {
					return errors.New("business closed")
				}
			case "Monday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliveryMonday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliveryMonday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeMonday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeMonday) || request.OrderDate.After(closingTimeMonday) {
					return errors.New("business closed")
				}
			case "Tuesday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliveryTuesday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliveryTuesday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeTuesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeTuesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeTuesday) || request.OrderDate.After(closingTimeTuesday) {
					return errors.New("business closed")
				}
			case "Wednesday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliveryWednesday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliveryWednesday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeWednesday) || request.OrderDate.After(closingTimeWednesday) {
					return errors.New("business closed")
				}
			case "Thursday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliveryThursday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliveryThursday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeThursday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeThursday) || request.OrderDate.After(closingTimeThursday) {
					return errors.New("business closed")
				}
			case "Friday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliveryFriday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliveryFriday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeFriday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeFriday) || request.OrderDate.After(closingTimeFriday) {
					return errors.New("business closed")
				}
			case "Saturday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeDeliverySaturday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeDeliverySaturday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeSaturday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeSaturday) || request.OrderDate.After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			}
		} else if request.OrderType == "OrderTypePickUp" {
			switch weekday {
			case "Sunday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeSunday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeSunday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeSunday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSunday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeSunday) || request.OrderDate.After(closingTimeSunday) {
					return errors.New("business closed")
				}
			case "Monday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeMonday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeMonday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeMonday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeMonday) || request.OrderDate.After(closingTimeMonday) {
					return errors.New("business closed")
				}
			case "Tuesday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeTuesday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeTuesday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeTuesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeTuesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeTuesday) || request.OrderDate.After(closingTimeTuesday) {
					return errors.New("business closed")
				}
			case "Wednesday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeWednesday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeWednesday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeWednesday) || request.OrderDate.After(closingTimeWednesday) {
					return errors.New("business closed")
				}
			case "Thursday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeThursday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeThursday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeThursday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeThursday) || request.OrderDate.After(closingTimeThursday) {
					return errors.New("business closed")
				}
			case "Friday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeFriday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeFriday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeFriday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeFriday) || request.OrderDate.After(closingTimeFriday) {
					return errors.New("business closed")
				}
			case "Saturday":
				splitOpening := strings.Split(businessScheduleRes.OpeningTimeSaturday, ":")
				splitClosing := strings.Split(businessScheduleRes.ClosingTimeSaturday, ":")
				openingHour, openingHourErr := strconv.Atoi(splitOpening[0])
				if openingHourErr != nil {
					return openingHourErr
				}
				openingMinutes, openingMinutesErr := strconv.Atoi(splitOpening[1])
				if openingMinutesErr != nil {
					return openingMinutesErr
				}
				closingHour, closingHourErr := strconv.Atoi(splitClosing[0])
				if closingHourErr != nil {
					return closingHourErr
				}
				closingMinutes, closingMinutesErr := strconv.Atoi(splitClosing[1])
				if closingMinutesErr != nil {
					return closingMinutesErr
				}
				openingTimeSaturday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if request.OrderDate.Before(openingTimeSaturday) || request.OrderDate.After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			}
		}
		_, createOrderedItemsErr := i.dao.NewOrderedRepository().BatchCreateOrderedItem(tx, &orderedItems)
		if createOrderedItemsErr != nil {
			return createOrderedItemsErr
		}
		createOrderRes, createOrderErr := i.dao.NewOrderRepository().CreateOrder(tx, &models.Order{Quantity: quantity, OrderType: request.OrderType, ResidenceType: request.ResidenceType, BuildingNumber: request.BuildingNumber, HouseNumber: request.HouseNumber, UserFk: authorizationTokenRes.UserFk, OrderDate: request.OrderDate, Coordinates: request.Coordinates, AuthorizationTokenFk: authorizationTokenRes.ID, BusinessFk: (*listCartItemRes)[0].BusinessFk, Price: price, CreateTime: createTime, UpdateTime: createTime})
		if createOrderErr != nil {
			return createOrderErr
		}
		_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: createOrderRes.Status, OrderFk: createOrderRes.ID})
		if createOrderLcErr != nil {
			return createOrderLcErr
		}
		unionOrderAndOrderedItems := make([]models.UnionOrderAndOrderedItem, 0, len(orderedItems))
		for _, item := range orderedItems {
			unionOrderAndOrderedItems = append(unionOrderAndOrderedItems, models.UnionOrderAndOrderedItem{OrderFk: createOrderRes.ID, OrderedItemFk: item.ID})
		}
		_, createUnionOrderAndOrderedItemsErr := i.dao.NewUnionOrderAndOrderedItemRepository().BatchCreateUnionOrderAndOrderedItem(tx, &unionOrderAndOrderedItems)
		if createUnionOrderAndOrderedItemsErr != nil {
			return createUnionOrderAndOrderedItemsErr
		}
		deleteCartItemErr := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{UserFk: authorizationTokenRes.UserFk})
		if deleteCartItemErr != nil {
			return deleteCartItemErr
		}
		response.Order = models.Order{Quantity: quantity, Status: createOrderRes.Status, OrderType: createOrderRes.OrderType, ResidenceType: createOrderRes.ResidenceType, BuildingNumber: createOrderRes.BuildingNumber, HouseNumber: createOrderRes.HouseNumber, BusinessFk: createOrderRes.BusinessFk, AuthorizationTokenFk: createOrderRes.AuthorizationTokenFk, UserFk: createOrderRes.UserFk, OrderDate: createOrderRes.OrderDate, Coordinates: createOrderRes.Coordinates, Price: price, CreateTime: createOrderRes.CreateTime, UpdateTime: createOrderRes.UpdateTime, ID: createOrderRes.ID}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) ListOrder(request *dto.ListOrderRequest) (*dto.ListOrderResponse, error) {
	var listOrderResponse dto.ListOrderResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		ordersRes, ordersErr := i.dao.NewOrderRepository().ListOrderWithBusiness(tx, &models.OrderBusiness{CreateTime: request.NextPage, UserFk: authorizationTokenRes.UserFk})
		if ordersErr != nil {
			return ordersErr
		}
		if len(*ordersRes) > 10 {
			*ordersRes = (*ordersRes)[:len(*ordersRes)-1]
			listOrderResponse.NextPage = (*ordersRes)[len(*ordersRes)-1].CreateTime
		} else if len(*ordersRes) > 0 {
			listOrderResponse.NextPage = (*ordersRes)[len(*ordersRes)-1].CreateTime
		} else {
			listOrderResponse.NextPage = request.NextPage
		}
		listOrderResponse.Orders = ordersRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &listOrderResponse, nil
}
