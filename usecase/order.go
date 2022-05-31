package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type OrderService interface {
	ListOrder(ctx context.Context, req *pb.ListOrderRequest, md *utils.ClientMetadata) (*pb.ListOrderResponse, error)
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: request.OrderId}, &[]string{})
		if unionOrderAndOrderedItemErr != nil {
			return unionOrderAndOrderedItemErr
		}
		orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
		for _, item := range *unionOrderAndOrderedItemRes {
			orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
		}
		orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks, &[]string{})
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		orderRes, orderErr := i.dao.NewOrderRepository().GetOrder(tx, &models.Order{ID: request.Id})
		if orderErr != nil {
			return orderErr
		}
		switch request.Status {
		case "OrderStatusTypeCanceled":
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: request.Id}, &[]string{})
			if unionOrderAndOrderedItemErr != nil {
				return unionOrderAndOrderedItemErr
			}
			orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
			for _, item := range *unionOrderAndOrderedItemRes {
				orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
			}
			orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks, &[]string{})
			if orderedItemsErr != nil {
				return orderedItemsErr
			}
			itemFks := make([]uuid.UUID, 0, len(*orderedItemsRes))
			for _, item := range *orderedItemsRes {
				itemFks = append(itemFks, *item.ItemId)
			}
			itemsRes, itemsErr := i.dao.NewItemQuery().ListItemInIds(tx, itemFks)
			if itemsErr != nil {
				return itemsErr
			}
			for _, item := range *orderedItemsRes {
				var index = -1
				for i, n := range *itemsRes {
					if *n.ID == *item.ItemId {
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
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: request.Status, OrderId: request.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime}
		case "OrderStatusTypeRejected":
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: request.Id}, &[]string{})
			if unionOrderAndOrderedItemErr != nil {
				return unionOrderAndOrderedItemErr
			}
			orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
			for _, item := range *unionOrderAndOrderedItemRes {
				orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
			}
			orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks, &[]string{})
			if orderedItemsErr != nil {
				return orderedItemsErr
			}
			itemFks := make([]uuid.UUID, 0, len(*orderedItemsRes))
			for _, item := range *orderedItemsRes {
				itemFks = append(itemFks, *item.ItemId)
			}
			itemsRes, itemsErr := i.dao.NewItemQuery().ListItemInIds(tx, itemFks)
			if itemsErr != nil {
				return itemsErr
			}
			for _, item := range *orderedItemsRes {
				var index = -1
				for i, n := range *itemsRes {
					if *n.ID == *item.ItemId {
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
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: request.Status, OrderId: request.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, ItemsQuantity: updateOrderRes.ItemsQuantity, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypePending":
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: request.Status, OrderId: request.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypeApproved":
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: request.Status, OrderId: request.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypeDone":
			if orderRes.Status != "OrderStatusTypeApproved" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: request.Status, OrderId: request.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypeReceived":
			if orderRes.Status != "OrderStatusTypeApproved" && orderRes.Status != "OrderStatusTypeDone" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: request.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: request.Status, OrderId: request.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		createTime := time.Now().UTC()
		weekday := request.OrderDate.Weekday().String()
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listCartItemRes, listCartItemErr := i.dao.NewCartItemRepository().ListCartItemInIds(tx, *request.CartItems, nil)
		if listCartItemErr != nil {
			return listCartItemErr
		}
		var price decimal.Decimal
		var quantity int32 = 0
		orderedItems := make([]models.OrderedItem, 0, len(*listCartItemRes))
		for _, item := range *listCartItemRes {
			itemPrice, itemPriceErr := decimal.NewFromString(item.Price)
			if itemPriceErr != nil {
				return itemPriceErr
			}
			price.Add(itemPrice)
			quantity += item.Quantity
			orderedItems = append(orderedItems, models.OrderedItem{Quantity: item.Quantity, Price: item.Price, CartItemId: item.ID, UserId: item.UserId, ItemId: item.ItemId})
		}
		businessScheduleRes, businessScheduleErr := i.dao.NewBusinessScheduleRepository().GetBusinessSchedule(tx, &models.BusinessSchedule{BusinessId: (*listCartItemRes)[0].BusinessId})
		if businessScheduleErr != nil {
			return businessScheduleErr
		}
		businessRes, businessErr := i.dao.NewBusinessQuery().GetBusinessWithLocation(tx, &models.Business{ID: (*listCartItemRes)[0].BusinessId, Coordinates: request.Coordinates})
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
				openingTimeSunday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeSunday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
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
				openingTimeTuesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeTuesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
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
				openingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
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
				openingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeWednesday := time.Date(request.OrderDate.Year(), request.OrderDate.Month(), request.OrderDate.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
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
		createOrderRes, createOrderErr := i.dao.NewOrderRepository().CreateOrder(tx, &models.Order{ItemsQuantity: quantity, OrderType: request.OrderType, ResidenceType: request.ResidenceType, UserId: authorizationTokenRes.UserId, OrderDate: request.OrderDate, Coordinates: request.Coordinates, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: (*listCartItemRes)[0].BusinessId, Price: price.String(), CreateTime: createTime, UpdateTime: createTime, Number: request.Number, Address: request.Address, Instructions: request.Instructions})
		if createOrderErr != nil {
			return createOrderErr
		}
		_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: createOrderRes.Status, OrderId: createOrderRes.ID})
		if createOrderLcErr != nil {
			return createOrderLcErr
		}
		unionOrderAndOrderedItems := make([]models.UnionOrderAndOrderedItem, 0, len(orderedItems))
		for _, item := range orderedItems {
			unionOrderAndOrderedItems = append(unionOrderAndOrderedItems, models.UnionOrderAndOrderedItem{OrderId: createOrderRes.ID, OrderedItemId: item.ID})
		}
		_, createUnionOrderAndOrderedItemsErr := i.dao.NewUnionOrderAndOrderedItemRepository().BatchCreateUnionOrderAndOrderedItem(tx, &unionOrderAndOrderedItems)
		if createUnionOrderAndOrderedItemsErr != nil {
			return createUnionOrderAndOrderedItemsErr
		}
		_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{UserId: authorizationTokenRes.UserId}, nil)
		if err != nil {
			return err
		}
		response.Order = models.Order{ItemsQuantity: quantity, Status: createOrderRes.Status, OrderType: createOrderRes.OrderType, ResidenceType: createOrderRes.ResidenceType, Number: createOrderRes.Number, BusinessId: createOrderRes.BusinessId, AuthorizationTokenId: createOrderRes.AuthorizationTokenId, UserId: createOrderRes.UserId, OrderDate: createOrderRes.OrderDate, Coordinates: createOrderRes.Coordinates, Price: price.String(), CreateTime: createOrderRes.CreateTime, UpdateTime: createOrderRes.UpdateTime, ID: createOrderRes.ID, Address: createOrderRes.Address, Instructions: createOrderRes.Instructions}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) ListOrder(ctx context.Context, req *pb.ListOrderRequest, md *utils.ClientMetadata) (*pb.ListOrderResponse, error) {
	var ordersRes *[]models.OrderBusiness
	var ordersErr error
	var nextPage time.Time
	if req.NextPage == nil {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	var res pb.ListOrderResponse
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		ordersRes, ordersErr = i.dao.NewOrderRepository().ListOrderWithBusiness(tx, &models.OrderBusiness{CreateTime: nextPage, UserId: authorizationTokenRes.UserId})
		if ordersErr != nil {
			return ordersErr
		}
		if len(*ordersRes) > 10 {
			*ordersRes = (*ordersRes)[:len(*ordersRes)-1]
			res.NextPage = timestamppb.New((*ordersRes)[len(*ordersRes)-1].CreateTime)
		} else if len(*ordersRes) > 0 {
			res.NextPage = timestamppb.New((*ordersRes)[len(*ordersRes)-1].CreateTime)
		} else {
			res.NextPage = timestamppb.New(nextPage)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	ordersResponse := make([]*pb.Order, 0, len(*ordersRes))
	for _, item := range *ordersRes {
		ordersResponse = append(ordersResponse, &pb.Order{
			Id:           item.ID.String(),
			BusinessName: item.BusinessName,
			Quantity:     item.Quantity,
			Price:        item.Price,
			Number:       item.Number, Address: item.Address,
			Instructions:  item.Instructions,
			UserId:        item.UserId.String(),
			OrderDate:     timestamppb.New(item.OrderDate),
			Status:        *utils.ParseOrderStatusType(&item.Status),
			OrderType:     *utils.ParseOrderType(&item.OrderType),
			ResidenceType: *utils.ParseOrderResidenceType(&item.ResidenceType),
			Coordinates:   &pb.Point{Latitude: item.Coordinates.Coords()[1], Longitude: item.Coordinates.Coords()[0]},
			BusinessId:    item.BusinessId.String(),
			CreateTime:    timestamppb.New(item.CreateTime),
			UpdateTime:    timestamppb.New(item.UpdateTime),
		})
	}
	res.Orders = ordersResponse
	return &res, nil
}
