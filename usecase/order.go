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
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type OrderService interface {
	ListOrder(ctx context.Context, req *pb.ListOrderRequest, md *utils.ClientMetadata) (*pb.ListOrderResponse, error)
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*dto.CreateOrderResponse, error)
	UpdateOrder(req *dto.UpdateOrderRequest) (*dto.UpdateOrderResponse, error)
	ListOrderedItemWithItem(req *dto.ListOrderedItemRequest) (*dto.ListOrderedItemResponse, error)
}

type orderService struct {
	dao repository.DAO
}

func NewOrderService(dao repository.DAO) OrderService {
	return &orderService{dao: dao}
}

func (i *orderService) ListOrderedItemWithItem(req *dto.ListOrderedItemRequest) (*dto.ListOrderedItemResponse, error) {
	var response dto.ListOrderedItemResponse
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &req.Metadata.Get("authorization")[0]}
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
		_, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: req.OrderId}, &[]string{})
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

func (i *orderService) UpdateOrder(req *dto.UpdateOrderRequest) (*dto.UpdateOrderResponse, error) {
	var response dto.UpdateOrderResponse
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &req.Metadata.Get("authorization")[0]}
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
		orderRes, orderErr := i.dao.NewOrderRepository().GetOrder(tx, &models.Order{ID: req.Id})
		if orderErr != nil {
			return orderErr
		}
		switch req.Status {
		case "OrderStatusTypeCanceled":
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: req.Id}, &[]string{})
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
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: req.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status, OrderId: req.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime}
		case "OrderStatusTypeRejected":
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: req.Id}, &[]string{})
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
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: req.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status, OrderId: req.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, ItemsQuantity: updateOrderRes.ItemsQuantity, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypePending":
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: req.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status, OrderId: req.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypeApproved":
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: req.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status, OrderId: req.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypeDone":
			if orderRes.Status != "OrderStatusTypeApproved" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: req.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status, OrderId: req.Id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, OrderType: updateOrderRes.OrderType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId, UserId: updateOrderRes.UserId, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenId: updateOrderRes.AuthorizationTokenId, OrderDate: updateOrderRes.OrderDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime, Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case "OrderStatusTypeReceived":
			if orderRes.Status != "OrderStatusTypeApproved" && orderRes.Status != "OrderStatusTypeDone" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: req.Id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status, OrderId: req.Id, CreateTime: updateOrderRes.UpdateTime})
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

func (i *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*dto.CreateOrderResponse, error) {
	var response dto.CreateOrderResponse
	cartItems := make([]uuid.UUID, 0, len(req.CartItems))
	for _, item := range req.CartItems {
		cartItems = append(cartItems, uuid.MustParse(item))
	}
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		createTime := time.Now().UTC()
		weekday := req.OrderDate.AsTime().Weekday().String()
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
		listCartItemRes, listCartItemErr := i.dao.NewCartItemRepository().ListCartItemInIds(tx, cartItems, nil)
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
		businessRes, businessErr := i.dao.NewBusinessQuery().GetBusinessWithLocation(tx, &models.Business{ID: (*listCartItemRes)[0].BusinessId, Coordinates: location})
		if businessErr != nil {
			return businessErr
		}
		previousTime := createTime
		previousTime = previousTime.AddDate(0, int(businessRes.TimeMarginOrderMonth), int(businessRes.TimeMarginOrderDay))
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderHour) * time.Hour)
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderMinute) * time.Minute)
		if req.OrderDate.AsTime().Before(previousTime) {
			return errors.New("invalid schedule")
		}
		if req.OrderType.String() == "OrderTypeHomeDelivery" {
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
				openingTimeSunday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeSunday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if req.OrderDate.AsTime().Before(openingTimeSunday) || req.OrderDate.AsTime().After(closingTimeSunday) {
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
				openingTimeMonday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeMonday) || req.OrderDate.AsTime().After(closingTimeMonday) {
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
				openingTimeTuesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeTuesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if req.OrderDate.AsTime().Before(openingTimeTuesday) || req.OrderDate.AsTime().After(closingTimeTuesday) {
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
				openingTimeWednesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeWednesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if req.OrderDate.AsTime().Before(openingTimeWednesday) || req.OrderDate.AsTime().After(closingTimeWednesday) {
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
				openingTimeThursday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeThursday) || req.OrderDate.AsTime().After(closingTimeThursday) {
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
				openingTimeFriday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeFriday) || req.OrderDate.AsTime().After(closingTimeFriday) {
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
				openingTimeSaturday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeSaturday) || req.OrderDate.AsTime().After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			}
		} else if req.OrderType.String() == "OrderTypePickUp" {
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
				openingTimeSunday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSunday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeSunday) || req.OrderDate.AsTime().After(closingTimeSunday) {
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
				openingTimeMonday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeMonday) || req.OrderDate.AsTime().After(closingTimeMonday) {
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
				openingTimeTuesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeTuesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeTuesday) || req.OrderDate.AsTime().After(closingTimeTuesday) {
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
				openingTimeWednesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeWednesday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if req.OrderDate.AsTime().Before(openingTimeWednesday) || req.OrderDate.AsTime().After(closingTimeWednesday) {
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
				openingTimeThursday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeThursday) || req.OrderDate.AsTime().After(closingTimeThursday) {
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
				openingTimeFriday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeFriday) || req.OrderDate.AsTime().After(closingTimeFriday) {
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
				openingTimeSaturday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(req.OrderDate.AsTime().Year(), req.OrderDate.AsTime().Month(), req.OrderDate.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if req.OrderDate.AsTime().Before(openingTimeSaturday) || req.OrderDate.AsTime().After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			}
		}
		_, createOrderedItemsErr := i.dao.NewOrderedRepository().BatchCreateOrderedItem(tx, &orderedItems)
		if createOrderedItemsErr != nil {
			return createOrderedItemsErr
		}
		createOrderRes, createOrderErr := i.dao.NewOrderRepository().CreateOrder(tx, &models.Order{ItemsQuantity: quantity, OrderType: req.OrderType.String(), ResidenceType: req.ResidenceType.String(), UserId: authorizationTokenRes.UserId, OrderDate: req.OrderDate.AsTime(), Coordinates: location, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: (*listCartItemRes)[0].BusinessId, Price: price.String(), CreateTime: createTime, UpdateTime: createTime, Number: req.Number, Address: req.Address, Instructions: req.Instructions})
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
