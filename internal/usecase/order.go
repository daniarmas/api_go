package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/internal/repository"
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
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error)
	UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error)
	ListOrderedItemWithItem(ctx context.Context, req *pb.ListOrderedItemRequest, md *utils.ClientMetadata) (*pb.ListOrderedItemResponse, error)
}

type orderService struct {
	dao   repository.Repository
	sqldb *sqldb.Sql
}

func NewOrderService(dao repository.Repository, sqldb *sqldb.Sql) OrderService {
	return &orderService{dao: dao, sqldb: sqldb}
}

func (i *orderService) ListOrderedItemWithItem(ctx context.Context, req *pb.ListOrderedItemRequest, md *utils.ClientMetadata) (*pb.ListOrderedItemResponse, error) {
	var res pb.ListOrderedItemResponse
	orderId := uuid.MustParse(req.OrderId)
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		_, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &entity.UnionOrderAndOrderedItem{OrderId: &orderId}, &[]string{"id", "order_id", "ordered_item_id"})
		if unionOrderAndOrderedItemErr != nil {
			return unionOrderAndOrderedItemErr
		}
		orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
		for _, item := range *unionOrderAndOrderedItemRes {
			orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
		}
		orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks, &[]string{"id", "name", "price_cup", "quantity", "item_id", "cart_item_id", "user_id", "create_time", "update_time"})
		if orderedItemsErr != nil {
			return orderedItemsErr
		}
		orderedItems := make([]*pb.OrderedItem, 0, len(*orderedItemsRes))
		for _, item := range *orderedItemsRes {
			orderedItems = append(orderedItems, &pb.OrderedItem{Id: item.ID.String(), Name: item.Name, PriceCup: item.PriceCup, ItemId: item.ItemId.String(), Quantity: item.Quantity, UserId: item.UserId.String(), CreateTime: timestamppb.New(item.CreateTime), UpdateTime: timestamppb.New(item.UpdateTime), CartItemId: item.CartItemId.String()})
		}
		res.OrderedItems = orderedItems
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *orderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error) {
	var res *pb.Order
	id := uuid.MustParse(req.Order.Id)
	var cancelReasons string
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		_, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		orderRes, orderErr := i.dao.NewOrderRepository().GetOrder(tx, &entity.Order{ID: &id})
		if orderErr != nil {
			return orderErr
		}
		switch req.Order.Status {
		case pb.OrderStatusType_OrderStatusTypePending:
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
		case pb.OrderStatusType_OrderStatusTypeRejected, pb.OrderStatusType_OrderStatusTypeCanceled:
			if req.Order.Status == pb.OrderStatusType_OrderStatusTypeCanceled && orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			} else if req.Order.Status == pb.OrderStatusType_OrderStatusTypeRejected && orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &entity.UnionOrderAndOrderedItem{OrderId: &id}, &[]string{"ordered_item_id"})
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
			itemsRes, itemsErr := i.dao.NewItemRepository().ListItemInIds(tx, itemFks, nil)
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
				_, updateItemsErr := i.dao.NewItemRepository().UpdateItem(tx, &entity.Item{ID: item.ID}, &item)
				if updateItemsErr != nil {
					return updateItemsErr
				}
			}
			cancelReasons = req.Order.CancelReasons
		case pb.OrderStatusType_OrderStatusTypeApproved:
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
		case pb.OrderStatusType_OrderStatusTypeDone:
			if orderRes.Status != "OrderStatusTypeApproved" {
				return errors.New("status error")
			}
		case pb.OrderStatusType_OrderStatusTypeReceived:
			if orderRes.Status != "OrderStatusTypeApproved" && orderRes.Status != "OrderStatusTypeDone" {
				return errors.New("status error")
			}
		}
		updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &entity.Order{ID: &id}, &entity.Order{Status: req.Order.Status.String(), CancelReasons: cancelReasons})
		if updateOrderErr != nil {
			return updateOrderErr
		}
		_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &entity.OrderLifecycle{Status: req.Order.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
		if createOrderLcErr != nil {
			return createOrderLcErr
		}
		res = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), PriceCup: updateOrderRes.PriceCup, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions, ShortId: updateOrderRes.ShortId, CancelReasons: updateOrderRes.CancelReasons, BusinessName: updateOrderRes.BusinessName, ItemsQuantity: updateOrderRes.ItemsQuantity}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error) {
	var res *pb.Order
	var cartItems []uuid.UUID
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		createTime := time.Now().UTC()
		orderTimeWeekday := req.OrderTime.AsTime()
		orderTimeWeekdayHour, _, _ := orderTimeWeekday.Clock()
		if orderTimeWeekdayHour >= 0 && orderTimeWeekdayHour <= 4 {
			orderTimeWeekday = orderTimeWeekday.AddDate(0, 0, 1)
		}
		weekday := orderTimeWeekday.Local().Weekday().String()
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		listCartItemRes, listCartItemErr := i.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, &[]string{"id", "business_id", "price_cup", "item_id", "user_id", "quantity", "name"})
		if listCartItemErr != nil {
			return listCartItemErr
		} else if listCartItemRes == nil {
			return errors.New("cart items not found")
		}
		for _, item := range *listCartItemRes {
			cartItems = append(cartItems, *item.ID)
		}
		var price_cup decimal.Decimal
		var quantity int32 = 0
		orderedItems := make([]entity.OrderedItem, 0, len(*listCartItemRes))
		for _, item := range *listCartItemRes {
			itemPriceCup, itemPriceCupErr := decimal.NewFromString(item.PriceCup)
			if itemPriceCupErr != nil {
				return itemPriceCupErr
			}
			price_cup = price_cup.Add(itemPriceCup)
			quantity += item.Quantity
			orderedItems = append(orderedItems, entity.OrderedItem{Quantity: item.Quantity, PriceCup: item.PriceCup, CartItemId: item.ID, UserId: item.UserId, ItemId: item.ItemId, Name: item.Name})
		}
		businessScheduleRes, businessScheduleErr := i.dao.NewBusinessScheduleRepository().GetBusinessSchedule(tx, &entity.BusinessSchedule{BusinessId: (*listCartItemRes)[0].BusinessId}, nil)
		if businessScheduleErr != nil {
			return businessScheduleErr
		}
		businessRes, businessErr := i.dao.NewBusinessRepository().GetBusinessWithDistance(tx, &entity.Business{ID: (*listCartItemRes)[0].BusinessId, Coordinates: location})
		if businessErr != nil {
			return businessErr
		}
		orderTime := req.OrderTime.AsTime()
		orderTimeHour, _, _ := orderTime.Clock()
		if orderTimeHour >= 0 && orderTimeHour <= 4 {
			orderTime = orderTime.AddDate(0, 0, 1)
		}
		previousTime := createTime
		previousTime = previousTime.AddDate(0, int(businessRes.TimeMarginOrderMonth), int(businessRes.TimeMarginOrderDay))
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderHour) * time.Hour)
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderMinute) * time.Minute)
		fmt.Printf("ordertime: %s\n", orderTime)
		fmt.Printf("previoustime: %s", previousTime)
		if orderTime.Before(previousTime) {
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
				openingTimeSunday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.Local)
				closingTimeSunday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.Local)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSunday = closingTimeSunday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeSunday) || orderTimeWeekday.After(closingTimeSunday) {
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
				openingTimeMonday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeMonday = closingTimeMonday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeMonday) || orderTimeWeekday.After(closingTimeMonday) {
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
				openingTimeTuesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeTuesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeTuesday = closingTimeTuesday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeTuesday) || orderTimeWeekday.After(closingTimeTuesday) {
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
				openingTimeWednesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeWednesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeWednesday = closingTimeWednesday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeWednesday) || orderTimeWeekday.After(closingTimeWednesday) {
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
				openingTimeThursday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeThursday = closingTimeThursday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeThursday) || orderTimeWeekday.After(closingTimeThursday) {
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
				openingTimeFriday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeFriday = closingTimeFriday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeFriday) || orderTimeWeekday.After(closingTimeFriday) {
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
				openingTimeSaturday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSaturday = closingTimeSaturday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeSaturday) || orderTimeWeekday.After(closingTimeSaturday) {
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
				openingTimeSunday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSunday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSunday = closingTimeSunday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeSunday) || orderTimeWeekday.After(closingTimeSunday) {
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
				openingTimeMonday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeMonday = closingTimeMonday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeMonday) || orderTimeWeekday.After(closingTimeMonday) {
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
				openingTimeTuesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeTuesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeTuesday = closingTimeTuesday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeTuesday) || orderTimeWeekday.After(closingTimeTuesday) {
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
				var openingTimeDay, closingTimeDay int
				openingTimeDay = time.Now().Day()
				if closingHour >= 0 && closingHour <= 4 {
					closingTimeDay = orderTimeWeekday.Add(time.Duration(24) * time.Hour).Day()
				}
				openingTimeWednesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), openingTimeDay, openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeWednesday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), closingTimeDay, closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeWednesday = closingTimeWednesday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeWednesday) || orderTimeWeekday.After(closingTimeWednesday) {
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
				openingTimeThursday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeThursday = closingTimeThursday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeThursday) || orderTimeWeekday.After(closingTimeThursday) {
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
				openingTimeFriday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeFriday = closingTimeFriday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeFriday) || orderTimeWeekday.After(closingTimeFriday) {
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
				openingTimeSaturday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(orderTimeWeekday.Year(), orderTimeWeekday.Month(), orderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSaturday = closingTimeSaturday.AddDate(0, 0, 1)
				}
				if orderTimeWeekday.Before(openingTimeSaturday) || orderTimeWeekday.After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			}
		}
		_, createOrderedItemsErr := i.dao.NewOrderedRepository().BatchCreateOrderedItem(tx, &orderedItems)
		if createOrderedItemsErr != nil {
			return createOrderedItemsErr
		}
		createOrderRes, createOrderErr := i.dao.NewOrderRepository().CreateOrder(tx, &entity.Order{ItemsQuantity: quantity, OrderType: req.OrderType.String(), UserId: authorizationTokenRes.UserId, OrderTime: req.OrderTime.AsTime().UTC(), Coordinates: location, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: (*listCartItemRes)[0].BusinessId, PriceCup: price_cup.String(), CreateTime: createTime, UpdateTime: createTime, Number: req.Number, Address: req.Address, Instructions: req.Instructions, BusinessName: businessRes.Name})
		if createOrderErr != nil {
			return createOrderErr
		}
		_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &entity.OrderLifecycle{Status: createOrderRes.Status, OrderId: createOrderRes.ID})
		if createOrderLcErr != nil {
			return createOrderLcErr
		}
		unionOrderAndOrderedItems := make([]entity.UnionOrderAndOrderedItem, 0, len(orderedItems))
		for _, item := range orderedItems {
			unionOrderAndOrderedItems = append(unionOrderAndOrderedItems, entity.UnionOrderAndOrderedItem{OrderId: createOrderRes.ID, OrderedItemId: item.ID})
		}
		_, createUnionOrderAndOrderedItemsErr := i.dao.NewUnionOrderAndOrderedItemRepository().BatchCreateUnionOrderAndOrderedItem(tx, &unionOrderAndOrderedItems)
		if createUnionOrderAndOrderedItemsErr != nil {
			return createUnionOrderAndOrderedItemsErr
		}
		_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, nil)
		if err != nil {
			return err
		}
		res = &pb.Order{BusinessName: businessRes.Name, ItemsQuantity: quantity, Status: *utils.ParseOrderStatusType(&createOrderRes.Status), OrderType: *utils.ParseOrderType(&createOrderRes.OrderType), Number: createOrderRes.Number, BusinessId: createOrderRes.BusinessId.String(), UserId: createOrderRes.UserId.String(), OrderTime: timestamppb.New(createOrderRes.OrderTime), Coordinates: &pb.Point{Latitude: createOrderRes.Coordinates.FlatCoords()[0], Longitude: createOrderRes.Coordinates.FlatCoords()[1]}, PriceCup: price_cup.String(), CreateTime: timestamppb.New(createOrderRes.CreateTime), UpdateTime: timestamppb.New(createOrderRes.UpdateTime), Address: createOrderRes.Address, Instructions: createOrderRes.Instructions, Id: createOrderRes.ID.String(), ShortId: createOrderRes.ShortId}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderService) ListOrder(ctx context.Context, req *pb.ListOrderRequest, md *utils.ClientMetadata) (*pb.ListOrderResponse, error) {
	var ordersRes *[]entity.OrderBusiness
	var ordersErr error
	var nextPage time.Time
	if req.NextPage == nil {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	var res pb.ListOrderResponse
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		ordersRes, ordersErr = i.dao.NewOrderRepository().ListOrderWithBusiness(tx, &entity.OrderBusiness{CreateTime: nextPage, UserId: authorizationTokenRes.UserId})
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
			Id:            item.ID.String(),
			ShortId:       item.ShortId,
			CancelReasons: item.CancelReasons,
			BusinessName:  item.BusinessName,
			ItemsQuantity: item.ItemsQuantity,
			PriceCup:      item.PriceCup,
			Number:        item.Number, Address: item.Address,
			Instructions: item.Instructions,
			UserId:       item.UserId.String(),
			OrderTime:    timestamppb.New(item.OrderTime),
			Status:       *utils.ParseOrderStatusType(&item.Status),
			OrderType:    *utils.ParseOrderType(&item.OrderType),
			Coordinates:  &pb.Point{Latitude: item.Coordinates.Coords()[1], Longitude: item.Coordinates.Coords()[0]},
			BusinessId:   item.BusinessId.String(),
			CreateTime:   timestamppb.New(item.CreateTime),
			UpdateTime:   timestamppb.New(item.UpdateTime),
		})
	}
	res.Orders = ordersResponse
	return &res, nil
}
