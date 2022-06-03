package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/daniarmas/api_go/datasource"
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
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*pb.CreateOrderResponse, error)
	UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest, md *utils.ClientMetadata) (*pb.UpdateOrderResponse, error)
	ListOrderedItemWithItem(ctx context.Context, req *pb.ListOrderedItemRequest, md *utils.ClientMetadata) (*pb.ListOrderedItemResponse, error)
}

type orderService struct {
	dao repository.DAO
}

func NewOrderService(dao repository.DAO) OrderService {
	return &orderService{dao: dao}
}

func (i *orderService) ListOrderedItemWithItem(ctx context.Context, req *pb.ListOrderedItemRequest, md *utils.ClientMetadata) (*pb.ListOrderedItemResponse, error) {
	var res pb.ListOrderedItemResponse
	orderId := uuid.MustParse(req.OrderId)
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
		_, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: &orderId}, &[]string{"id", "order_id", "ordered_item_id"})
		if unionOrderAndOrderedItemErr != nil {
			return unionOrderAndOrderedItemErr
		}
		orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
		for _, item := range *unionOrderAndOrderedItemRes {
			orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
		}
		orderedItemsRes, orderedItemsErr := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, &orderedItemFks, &[]string{"id", "name", "price", "quantity", "item_id", "cart_item_id", "user_id", "create_time", "update_time"})
		if orderedItemsErr != nil {
			return orderedItemsErr
		}
		orderedItems := make([]*pb.OrderedItem, 0, len(*orderedItemsRes))
		for _, item := range *orderedItemsRes {
			orderedItems = append(orderedItems, &pb.OrderedItem{Id: item.ID.String(), Name: item.Name, Price: item.Price, ItemId: item.ItemId.String(), Quantity: item.Quantity, UserId: item.UserId.String(), CreateTime: timestamppb.New(item.CreateTime), UpdateTime: timestamppb.New(item.UpdateTime), CartItemId: item.CartItemId.String()})
		}
		res.OrderedItems = orderedItems
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *orderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest, md *utils.ClientMetadata) (*pb.UpdateOrderResponse, error) {
	var response pb.UpdateOrderResponse
	id := uuid.MustParse(req.Id)
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
		orderRes, orderErr := i.dao.NewOrderRepository().GetOrder(tx, &models.Order{ID: &id})
		if orderErr != nil {
			return orderErr
		}
		switch req.Status {
		case pb.OrderStatusType_OrderStatusTypeCanceled:
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: &id}, &[]string{})
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
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: &id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status.String()})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.ResidenceType), Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case pb.OrderStatusType_OrderStatusTypeRejected:
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &models.UnionOrderAndOrderedItem{OrderId: &id}, &[]string{})
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
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: &id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status.String()})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.ResidenceType), Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case pb.OrderStatusType_OrderStatusTypePending:
			if orderRes.Status != "OrderStatusTypeStarted" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: &id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status.String()})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.ResidenceType), Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case pb.OrderStatusType_OrderStatusTypeApproved:
			if orderRes.Status != "OrderStatusTypePending" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: &id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status.String()})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.ResidenceType), Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case pb.OrderStatusType_OrderStatusTypeDone:
			if orderRes.Status != "OrderStatusTypeApproved" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: &id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status.String()})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.ResidenceType), Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		case pb.OrderStatusType_OrderStatusTypeReceived:
			if orderRes.Status != "OrderStatusTypeApproved" && orderRes.Status != "OrderStatusTypeDone" {
				return errors.New("status error")
			}
			updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: &id, UserId: authorizationTokenRes.UserId}, &models.Order{Status: req.Status.String()})
			if updateOrderErr != nil {
				return updateOrderErr
			}
			_, createOrderLcErr := i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &models.OrderLifecycle{Status: req.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
			if createOrderLcErr != nil {
				return createOrderLcErr
			}
			response.Order = &pb.Order{Id: updateOrderRes.ID.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.ResidenceType), Price: updateOrderRes.Price, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, OrderTime: timestamppb.New(updateOrderRes.OrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*pb.CreateOrderResponse, error) {
	var response pb.CreateOrderResponse
	cartItems := make([]uuid.UUID, 0, len(req.CartItems))
	for _, item := range req.CartItems {
		cartItems = append(cartItems, uuid.MustParse(item))
	}
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		createTime := time.Now().UTC()
		weekday := req.OrderTime.AsTime().Local().Weekday().String()
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
		listCartItemRes, listCartItemErr := i.dao.NewCartItemRepository().ListCartItemInIds(tx, cartItems, &[]string{"id", "item_id", "user_id", "price", "quantity", "business_id", "name"})
		if listCartItemErr != nil {
			return listCartItemErr
		} else if listCartItemRes == nil {
			return errors.New("cart items not found")
		}
		var price decimal.Decimal
		var quantity int32 = 0
		orderedItems := make([]models.OrderedItem, 0, len(*listCartItemRes))
		for _, item := range *listCartItemRes {
			itemPrice, itemPriceErr := decimal.NewFromString(item.Price)
			if itemPriceErr != nil {
				return itemPriceErr
			}
			price = price.Add(itemPrice)
			quantity += item.Quantity
			orderedItems = append(orderedItems, models.OrderedItem{Quantity: item.Quantity, Price: item.Price, CartItemId: item.ID, UserId: item.UserId, ItemId: item.ItemId, Name: item.Name})
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
		if req.OrderTime.AsTime().UTC().Before(previousTime) {
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
				openingTimeSunday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local)
				closingTimeSunday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSunday = closingTimeSunday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeSunday) || req.OrderTime.AsTime().After(closingTimeSunday) {
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
				openingTimeMonday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeMonday = closingTimeMonday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeMonday) || req.OrderTime.AsTime().After(closingTimeMonday) {
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
				openingTimeTuesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeTuesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeTuesday = closingTimeTuesday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeTuesday) || req.OrderTime.AsTime().After(closingTimeTuesday) {
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
				openingTimeWednesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.Local).UTC()
				closingTimeWednesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.Local).UTC()
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeWednesday = closingTimeWednesday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeWednesday) || req.OrderTime.AsTime().After(closingTimeWednesday) {
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
				openingTimeThursday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeThursday = closingTimeThursday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeThursday) || req.OrderTime.AsTime().After(closingTimeThursday) {
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
				openingTimeFriday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeFriday = closingTimeFriday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeFriday) || req.OrderTime.AsTime().After(closingTimeFriday) {
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
				openingTimeSaturday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSaturday = closingTimeSaturday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeSaturday) || req.OrderTime.AsTime().After(closingTimeSaturday) {
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
				openingTimeSunday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSunday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSunday = closingTimeSunday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeSunday) || req.OrderTime.AsTime().After(closingTimeSunday) {
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
				openingTimeMonday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeMonday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeMonday = closingTimeMonday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeMonday) || req.OrderTime.AsTime().After(closingTimeMonday) {
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
				openingTimeTuesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeTuesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeTuesday = closingTimeTuesday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeTuesday) || req.OrderTime.AsTime().After(closingTimeTuesday) {
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
					closingTimeDay = req.OrderTime.AsTime().Add(time.Duration(24) * time.Hour).Day()
				}
				openingTimeWednesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), openingTimeDay, openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeWednesday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), closingTimeDay, closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeWednesday = closingTimeWednesday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeWednesday) || req.OrderTime.AsTime().After(closingTimeWednesday) {
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
				openingTimeThursday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeThursday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeThursday = closingTimeThursday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeThursday) || req.OrderTime.AsTime().After(closingTimeThursday) {
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
				openingTimeFriday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeFriday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeFriday = closingTimeFriday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeFriday) || req.OrderTime.AsTime().After(closingTimeFriday) {
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
				openingTimeSaturday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				closingTimeSaturday := time.Date(req.OrderTime.AsTime().Year(), req.OrderTime.AsTime().Month(), req.OrderTime.AsTime().Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					closingTimeSaturday = closingTimeSaturday.AddDate(0, 0, 1)
				}
				if req.OrderTime.AsTime().Before(openingTimeSaturday) || req.OrderTime.AsTime().After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			}
		}
		_, createOrderedItemsErr := i.dao.NewOrderedRepository().BatchCreateOrderedItem(tx, &orderedItems)
		if createOrderedItemsErr != nil {
			return createOrderedItemsErr
		}
		createOrderRes, createOrderErr := i.dao.NewOrderRepository().CreateOrder(tx, &models.Order{ItemsQuantity: quantity, OrderType: req.OrderType.String(), ResidenceType: req.ResidenceType.String(), UserId: authorizationTokenRes.UserId, OrderTime: req.OrderTime.AsTime().UTC(), Coordinates: location, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: (*listCartItemRes)[0].BusinessId, Price: price.String(), CreateTime: createTime, UpdateTime: createTime, Number: req.Number, Address: req.Address, Instructions: req.Instructions, BusinessName: businessRes.Name})
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
		response.Order = &pb.Order{BusinessName: businessRes.Name, ItemsQuantity: quantity, Status: *utils.ParseOrderStatusType(&createOrderRes.Status), OrderType: *utils.ParseOrderType(&createOrderRes.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&createOrderRes.ResidenceType), Number: createOrderRes.Number, BusinessId: createOrderRes.BusinessId.String(), UserId: createOrderRes.UserId.String(), OrderTime: timestamppb.New(createOrderRes.OrderTime), Coordinates: &pb.Point{Latitude: createOrderRes.Coordinates.FlatCoords()[0], Longitude: createOrderRes.Coordinates.FlatCoords()[1]}, Price: price.String(), CreateTime: timestamppb.New(createOrderRes.CreateTime), UpdateTime: timestamppb.New(createOrderRes.UpdateTime), Address: createOrderRes.Address, Instructions: createOrderRes.Instructions, Id: createOrderRes.ID.String(), ShortId: createOrderRes.ShortId}
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
			Id:            item.ID.String(),
			ShortId:       item.ShortId,
			CancelReasons: item.CancelReasons,
			BusinessName:  item.BusinessName,
			ItemsQuantity: item.ItemsQuantity,
			Price:         item.Price,
			Number:        item.Number, Address: item.Address,
			Instructions:  item.Instructions,
			UserId:        item.UserId.String(),
			OrderTime:     timestamppb.New(item.OrderTime),
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
