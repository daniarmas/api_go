package usecase

import (
	"context"
	"errors"

	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type OrderService interface {
	GetCheckoutInfo(ctx context.Context, req *pb.GetCheckoutInfoRequest, md *utils.ClientMetadata) (*pb.GetCheckoutInfoResponse, error)
	GetOrder(ctx context.Context, req *pb.GetOrderRequest, md *utils.ClientMetadata) (*pb.Order, error)
	ListOrder(ctx context.Context, req *pb.ListOrderRequest, md *utils.ClientMetadata) (*pb.ListOrderResponse, error)
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error)
	UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error)
	ListOrderedItemWithItem(ctx context.Context, req *pb.ListOrderedItemRequest, md *utils.ClientMetadata) (*pb.ListOrderedItemResponse, error)
}

type orderService struct {
	dao    repository.Repository
	config *config.Config
	sqldb  *sqldb.Sql
}

func NewOrderService(dao repository.Repository, sqldb *sqldb.Sql, config *config.Config) OrderService {
	return &orderService{dao: dao, sqldb: sqldb, config: config}
}

func (i *orderService) GetCheckoutInfo(ctx context.Context, req *pb.GetCheckoutInfoRequest, md *utils.ClientMetadata) (*pb.GetCheckoutInfoResponse, error) {
	var res pb.GetCheckoutInfoResponse
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		_, err = i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		businessId := uuid.MustParse(req.BusinessId)
		businessRes, err := i.dao.NewBusinessRepository().GetBusiness(tx, &entity.Business{ID: &businessId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("business not found")
		} else if err != nil {
			return err
		}
		businessIsInRange, err := i.dao.NewBusinessRepository().BusinessIsInRange(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Coordinates.Latitude, req.Coordinates.Longitude}).SetSRID(4326)}, businessRes.ID)
		if err != nil {
			return err
		}
		businessDistance, err := i.dao.NewBusinessRepository().GetBusinessDistance(tx, &entity.Business{ID: businessRes.ID}, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Coordinates.Latitude, req.Coordinates.Longitude}).SetSRID(4326)})
		if err != nil {
			return err
		}
		schedule, err := i.dao.NewBusinessScheduleRepository().GetBusinessSchedule(tx, &entity.BusinessSchedule{BusinessId: businessRes.ID})
		if err != nil {
			return err
		}
		listBusinessPaymentMethodsResult, err := i.dao.NewBusinessPaymentMethodRepository().ListBusinessPaymentMethodWithEnabled(ctx, tx, &entity.BusinessPaymentMethod{BusinessId: &businessId})
		if err != nil {
			return err
		}
		businessPaymentMethods := make([]*pb.BusinessPaymentMethod, 0, len(*listBusinessPaymentMethodsResult))
		for _, item := range *listBusinessPaymentMethodsResult {
			businessPaymentMethods = append(businessPaymentMethods, &pb.BusinessPaymentMethod{
				Id:              item.ID.String(),
				Type:            *utils.ParsePaymentMethodType(&item.Type),
				Address:         item.Address,
				Enabled:         item.Enabled,
				BusinessId:      item.BusinessId.String(),
				PaymentMethodId: item.PaymentMethodId.String(),
				CreateTime:      timestamppb.New(item.CreateTime),
				UpdateTime:      timestamppb.New(item.UpdateTime),
			})
		}
		res = pb.GetCheckoutInfoResponse{
			ServerTimeNow:          timestamppb.New(time.Now().UTC()),
			BusinessPaymentMethods: businessPaymentMethods,
			BusinessAddress:        businessRes.Address,
			Delivery:               businessRes.HomeDelivery,
			PickUp:                 businessRes.ToPickUp,
			DeliveryPriceCup:       businessRes.DeliveryPriceCup,
			BusinessCoordinates:    &pb.Point{Latitude: businessRes.Coordinates.Coords()[1], Longitude: businessRes.Coordinates.Coords()[0]},
			TimeMarginOrderMonth:   businessRes.TimeMarginOrderMonth,
			TimeMarginOrderDay:     businessRes.TimeMarginOrderDay,
			TimeMarginOrderHour:    businessRes.TimeMarginOrderHour,
			TimeMarginOrderMinute:  businessRes.TimeMarginOrderMinute,
			IsInRange:              *businessIsInRange,
			Distance:               businessDistance.Distance,
			BusinessSchedule: &pb.BusinessSchedule{
				Id:                         schedule.ID.String(),
				FirstOpeningTimeSunday:     timestamppb.New(schedule.FirstOpeningTimeSunday),
				FirstClosingTimeSunday:     timestamppb.New(schedule.FirstClosingTimeSunday),
				FirstOpeningTimeMonday:     timestamppb.New(schedule.FirstOpeningTimeMonday),
				FirstClosingTimeMonday:     timestamppb.New(schedule.FirstClosingTimeMonday),
				FirstOpeningTimeTuesday:    timestamppb.New(schedule.FirstOpeningTimeTuesday),
				FirstClosingTimeTuesday:    timestamppb.New(schedule.FirstClosingTimeTuesday),
				FirstOpeningTimeWednesday:  timestamppb.New(schedule.FirstOpeningTimeWednesday),
				FirstClosingTimeWednesday:  timestamppb.New(schedule.FirstClosingTimeWednesday),
				FirstOpeningTimeThursday:   timestamppb.New(schedule.FirstOpeningTimeThursday),
				FirstClosingTimeThursday:   timestamppb.New(schedule.FirstClosingTimeThursday),
				FirstOpeningTimeFriday:     timestamppb.New(schedule.FirstOpeningTimeFriday),
				FirstClosingTimeFriday:     timestamppb.New(schedule.FirstClosingTimeFriday),
				FirstOpeningTimeSaturday:   timestamppb.New(schedule.FirstOpeningTimeSaturday),
				FirstClosingTimeSaturday:   timestamppb.New(schedule.FirstClosingTimeSaturday),
				SecondOpeningTimeSunday:    timestamppb.New(schedule.SecondOpeningTimeSunday),
				SecondClosingTimeSunday:    timestamppb.New(schedule.SecondClosingTimeSunday),
				SecondOpeningTimeMonday:    timestamppb.New(schedule.SecondOpeningTimeMonday),
				SecondClosingTimeMonday:    timestamppb.New(schedule.SecondClosingTimeMonday),
				SecondOpeningTimeTuesday:   timestamppb.New(schedule.SecondOpeningTimeTuesday),
				SecondClosingTimeTuesday:   timestamppb.New(schedule.SecondClosingTimeTuesday),
				SecondOpeningTimeWednesday: timestamppb.New(schedule.SecondOpeningTimeWednesday),
				SecondClosingTimeWednesday: timestamppb.New(schedule.SecondClosingTimeWednesday),
				SecondOpeningTimeThursday:  timestamppb.New(schedule.SecondOpeningTimeThursday),
				SecondClosingTimeThursday:  timestamppb.New(schedule.SecondClosingTimeThursday),
				SecondOpeningTimeFriday:    timestamppb.New(schedule.SecondOpeningTimeFriday),
				SecondClosingTimeFriday:    timestamppb.New(schedule.SecondClosingTimeFriday),
				SecondOpeningTimeSaturday:  timestamppb.New(schedule.SecondOpeningTimeSaturday),
				SecondClosingTimeSaturday:  timestamppb.New(schedule.SecondClosingTimeSaturday),
			},
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *orderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest, md *utils.ClientMetadata) (*pb.Order, error) {
	var res pb.Order
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		id := uuid.MustParse(req.Id)
		order, err := i.dao.NewOrderRepository().GetOrder(tx, &entity.Order{ID: &id, UserId: authorizationTokenRes.UserId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("order not found")
		} else if err != nil {
			return err
		}
		unionOrderedItems, err := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &entity.UnionOrderAndOrderedItem{OrderId: order.ID})
		if err != nil {
			return err
		}
		orderedItemIds := make([]uuid.UUID, 0, len(*unionOrderedItems))
		for _, item := range *unionOrderedItems {
			orderedItemIds = append(orderedItemIds, *item.OrderedItemId)
		}
		orderedItemsRes, err := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, orderedItemIds)
		if err != nil {
			return err
		}
		orderedItems := make([]*pb.OrderedItem, 0, len(*orderedItemsRes))
		for _, item := range *orderedItemsRes {
			orderedItems = append(orderedItems, &pb.OrderedItem{
				Id:         item.ID.String(),
				Name:       item.Name,
				PriceCup:   item.PriceCup,
				CartItemId: item.CartItemId.String(),
				ItemId:     item.ItemId.String(),
				Quantity:   item.Quantity,
				UserId:     item.UserId.String(),
				CreateTime: timestamppb.New(item.CreateTime),
				UpdateTime: timestamppb.New(item.UpdateTime),
			})
		}
		res = pb.Order{
			Id:                order.ID.String(),
			BusinessName:      order.BusinessName,
			Status:            *utils.ParseOrderStatusType(&order.Status),
			OrderType:         *utils.ParseOrderType(&order.OrderType),
			Coordinates:       &pb.Point{Latitude: order.Coordinates.Coords()[0], Longitude: order.Coordinates.Coords()[1]},
			ItemsQuantity:     order.ItemsQuantity,
			ShortId:           order.ShortId,
			Number:            order.Number,
			Address:           order.Address,
			Instructions:      order.Instructions,
			CancelReasons:     order.CancelReasons,
			PriceCup:          order.PriceCup,
			BusinessId:        order.BusinessId.String(),
			BusinessThumbnail: i.config.BusinessAvatarBulkName + "/" + order.BusinessThumbnail,
			UserId:            order.UserId.String(),
			OrderedItems:      orderedItems,
			CreateTime:        timestamppb.New(order.CreateTime),
			UpdateTime:        timestamppb.New(order.UpdateTime),
			DeliveryPriceCup:  order.DeliveryPriceCup,
			StartOrderTime:    timestamppb.New(order.StartOrderTime),
			EndOrderTime:      timestamppb.New(order.EndOrderTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *orderService) ListOrderedItemWithItem(ctx context.Context, req *pb.ListOrderedItemRequest, md *utils.ClientMetadata) (*pb.ListOrderedItemResponse, error) {
	var res pb.ListOrderedItemResponse
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		_, err = i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		orderId := uuid.MustParse(req.OrderId)
		unionOrderAndOrderedItemRes, err := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &entity.UnionOrderAndOrderedItem{OrderId: &orderId})
		if err != nil {
			return err
		}
		orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
		for _, item := range *unionOrderAndOrderedItemRes {
			orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
		}
		orderedItemsRes, err := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, orderedItemFks)
		if err != nil {
			return err
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
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		_, err = i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		orderRes, err := i.dao.NewOrderRepository().GetOrder(tx, &entity.Order{ID: &id})
		if err != nil && err.Error() == "record not found" {
			return errors.New("order not found")
		} else if err != nil {
			return err
		}
		switch req.Order.Status {
		case pb.OrderStatusType_OrderStatusTypeExpired:
			return errors.New("status error")
		case pb.OrderStatusType_OrderStatusTypeRejected:
			if orderRes.Status != "OrderStatusTypeOrdered" {
				return errors.New("status error")
			}
			unionOrderAndOrderedItemRes, unionOrderAndOrderedItemErr := i.dao.NewUnionOrderAndOrderedItemRepository().ListUnionOrderAndOrderedItem(tx, &entity.UnionOrderAndOrderedItem{OrderId: &id})
			if unionOrderAndOrderedItemErr != nil {
				return unionOrderAndOrderedItemErr
			}
			orderedItemFks := make([]uuid.UUID, 0, len(*unionOrderAndOrderedItemRes))
			for _, item := range *unionOrderAndOrderedItemRes {
				orderedItemFks = append(orderedItemFks, *item.OrderedItemId)
			}
			orderedItemsRes, err := i.dao.NewOrderedRepository().ListOrderedItemByIds(tx, orderedItemFks)
			if err != nil {
				return err
			}
			itemFks := make([]uuid.UUID, 0, len(*orderedItemsRes))
			for _, item := range *orderedItemsRes {
				itemFks = append(itemFks, *item.ItemId)
			}
			itemsRes, err := i.dao.NewItemRepository().ListItemInIds(ctx, tx, itemFks)
			if err != nil {
				return err
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
				_, err := i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.ItemBusiness{ID: item.ID}, &item)
				if err != nil {
					return err
				}
			}
			cancelReasons = req.Order.CancelReasons
			// case pb.OrderStatusType_OrderStatusTypeApproved:
			// 	if orderRes.Status != "OrderStatusTypePending" {
			// 		return errors.New("status error")
			// 	}
			// case pb.OrderStatusType_OrderStatusTypeDone:
			// 	if orderRes.Status != "OrderStatusTypeApproved" {
			// 		return errors.New("status error")
			// 	}
			// case pb.OrderStatusType_OrderStatusTypeReceived:
			// 	if orderRes.Status != "OrderStatusTypeApproved" && orderRes.Status != "OrderStatusTypeDone" {
			// 		return errors.New("status error")
			// 	}
		}
		updateOrderRes, err := i.dao.NewOrderRepository().UpdateOrder(tx, &entity.Order{ID: &id}, &entity.Order{Status: req.Order.Status.String(), CancelReasons: cancelReasons})
		if err != nil {
			return err
		}
		_, err = i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &entity.OrderLifecycle{Status: req.Order.Status.String(), OrderId: &id, CreateTime: updateOrderRes.UpdateTime})
		if err != nil {
			return err
		}
		res = &pb.Order{Id: updateOrderRes.ID.String(), DeliveryPriceCup: updateOrderRes.DeliveryPriceCup, BusinessThumbnail: i.config.BusinessAvatarBulkName + "/" + updateOrderRes.BusinessThumbnail, Status: *utils.ParseOrderStatusType(&updateOrderRes.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.OrderType), PriceCup: updateOrderRes.PriceCup, BusinessId: updateOrderRes.BusinessId.String(), UserId: updateOrderRes.UserId.String(), Coordinates: &pb.Point{Latitude: updateOrderRes.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Coordinates.FlatCoords()[1]}, StartOrderTime: timestamppb.New(updateOrderRes.StartOrderTime), EndOrderTime: timestamppb.New(updateOrderRes.EndOrderTime), CreateTime: timestamppb.New(updateOrderRes.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.UpdateTime), Number: updateOrderRes.Number, Address: updateOrderRes.Address, Instructions: updateOrderRes.Instructions, ShortId: updateOrderRes.ShortId, CancelReasons: updateOrderRes.CancelReasons, BusinessName: updateOrderRes.BusinessName, ItemsQuantity: updateOrderRes.ItemsQuantity}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, md *utils.ClientMetadata) (*pb.Order, error) {
	var res *pb.Order
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		startOrderTimeWeekday := req.StartOrderTime.AsTime()
		endOrderTimeWeekday := req.EndOrderTime.AsTime()
		startOrderTimeWeekdayHour, _, _ := startOrderTimeWeekday.Clock()
		if startOrderTimeWeekdayHour >= 0 && startOrderTimeWeekdayHour <= 4 {
			startOrderTimeWeekday = startOrderTimeWeekday.AddDate(0, 0, 1)
		}
		weekday := startOrderTimeWeekday.Local().Weekday().String()
		listCartItemRes, err := i.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId})
		if err != nil {
			return err
		} else if listCartItemRes == nil || len(*listCartItemRes) == 0 {
			return errors.New("cart items not found")
		}
		var cartItems []uuid.UUID
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
		businessScheduleRes, err := i.dao.NewBusinessScheduleRepository().GetBusinessSchedule(tx, &entity.BusinessSchedule{BusinessId: (*listCartItemRes)[0].BusinessId})
		if err != nil {
			return err
		}
		userAddressId := uuid.MustParse(req.UserAddressId)
		userAddress, err := i.dao.NewUserAddressRepository().GetUserAddress(tx, &entity.UserAddress{ID: &userAddressId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("user address not found")
		} else if err != nil {
			return err
		}
		businessRes, err := i.dao.NewBusinessRepository().GetBusinessWithDistance(tx, &entity.Business{ID: (*listCartItemRes)[0].BusinessId}, userAddress.Coordinates)
		if err != nil {
			return err
		}
		// Check if the time when the order is placed is between the choseen interval.
		createTime := time.Now().UTC()
		previousTime := createTime
		previousTime = previousTime.AddDate(0, int(businessRes.TimeMarginOrderMonth), int(businessRes.TimeMarginOrderDay))
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderHour) * time.Hour)
		previousTime = previousTime.Add(time.Duration(businessRes.TimeMarginOrderMinute) * time.Minute)
		// Add 10 minutes for payments.
		previousTime = previousTime.Add(10 * time.Minute)
		if previousTime.After(req.EndOrderTime.AsTime()) {
			return errors.New("not fulfilled the previous time of the business")
		}
		// If the order is for delivery, check if the location is in the delivery range of the business
		location := userAddress.Coordinates
		// if req.OrderType == pb.OrderType_OrderTypeHomeDelivery {
		// 	isInRange, err := i.dao.NewBusinessRepository().BusinessIsInRange(tx, location, businessRes.ID)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if !*isInRange {
		// 		return errors.New("business not in range")
		// 	}
		// }
		var zeroTime time.Time
		switch weekday {
		case "Sunday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeSunday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeSunday.Clock()
			openingTimeSunday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeSunday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeSunday = closingTimeSunday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeSunday == zeroTime {
				if startOrderTimeWeekday.Before(openingTimeSunday) || endOrderTimeWeekday.After(closingTimeSunday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeSunday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeSunday.Clock()
				secondOpeningTimeSunday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeSunday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeSunday = secondClosingTimeSunday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeSunday) || (endOrderTimeWeekday.After(closingTimeSunday) && startOrderTimeWeekday.Before(secondOpeningTimeSunday))) || (startOrderTimeWeekday.After(secondClosingTimeSunday) || endOrderTimeWeekday.After(secondClosingTimeSunday)) {
					return errors.New("business closed")
				}
			}
		case "Monday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeMonday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeMonday.Clock()
			openingTimeMonday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeMonday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeMonday = closingTimeMonday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeMonday == zeroTime {
				if startOrderTimeWeekday.Before(openingTimeMonday) || endOrderTimeWeekday.After(closingTimeMonday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeMonday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeMonday.Clock()
				secondOpeningTimeMonday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeMonday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeMonday = secondClosingTimeMonday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeMonday) || (endOrderTimeWeekday.After(closingTimeMonday) && startOrderTimeWeekday.Before(secondOpeningTimeMonday))) || (startOrderTimeWeekday.After(secondClosingTimeMonday) || endOrderTimeWeekday.After(secondClosingTimeMonday)) {
					return errors.New("business closed")
				}
			}
		case "Tuesday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeTuesday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeTuesday.Clock()
			openingTimeTuesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeTuesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeTuesday = closingTimeTuesday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeTuesday.String() == "" {
				if startOrderTimeWeekday.Before(openingTimeTuesday) || endOrderTimeWeekday.After(closingTimeTuesday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeTuesday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeTuesday.Clock()
				secondOpeningTimeTuesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeTuesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeTuesday = secondClosingTimeTuesday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeTuesday) || (endOrderTimeWeekday.After(closingTimeTuesday) && startOrderTimeWeekday.Before(secondOpeningTimeTuesday))) || (startOrderTimeWeekday.After(secondClosingTimeTuesday) || endOrderTimeWeekday.After(secondClosingTimeTuesday)) {
					return errors.New("business closed")
				}
			}
		case "Wednesday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeWednesday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeWednesday.Clock()
			openingTimeWednesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeWednesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeWednesday = closingTimeWednesday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeWednesday == zeroTime {
				if startOrderTimeWeekday.Before(openingTimeWednesday) || endOrderTimeWeekday.After(closingTimeWednesday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeWednesday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeWednesday.Clock()
				secondOpeningTimeWednesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeWednesday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeWednesday = secondClosingTimeWednesday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeWednesday) || (endOrderTimeWeekday.After(closingTimeWednesday) && startOrderTimeWeekday.Before(secondOpeningTimeWednesday))) || (startOrderTimeWeekday.After(secondClosingTimeWednesday) || endOrderTimeWeekday.After(secondClosingTimeWednesday)) {
					return errors.New("business closed")
				}
			}
		case "Thursday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeThursday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeThursday.Clock()
			openingTimeThursday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeThursday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeThursday = closingTimeThursday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeThursday == zeroTime {
				if startOrderTimeWeekday.Before(openingTimeThursday) || endOrderTimeWeekday.After(closingTimeThursday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeThursday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeThursday.Clock()
				secondOpeningTimeThursday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeThursday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeThursday = secondClosingTimeThursday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeThursday) || (endOrderTimeWeekday.After(closingTimeThursday) && startOrderTimeWeekday.Before(secondOpeningTimeThursday))) || (startOrderTimeWeekday.After(secondClosingTimeThursday) || endOrderTimeWeekday.After(secondClosingTimeThursday)) {
					return errors.New("business closed")
				}
			}
		case "Friday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeFriday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeFriday.Clock()
			openingTimeFriday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeFriday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeFriday = closingTimeFriday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeFriday == zeroTime {
				if startOrderTimeWeekday.Before(openingTimeFriday) || endOrderTimeWeekday.After(closingTimeFriday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.SecondOpeningTimeFriday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.SecondClosingTimeFriday.Clock()
				secondOpeningTimeFriday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeFriday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeFriday = secondClosingTimeFriday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeFriday) || (endOrderTimeWeekday.After(closingTimeFriday) && startOrderTimeWeekday.Before(secondOpeningTimeFriday))) || (startOrderTimeWeekday.After(secondClosingTimeFriday) || endOrderTimeWeekday.After(secondClosingTimeFriday)) {
					return errors.New("business closed")
				}
			}
		case "Saturday":
			openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeSaturday.Clock()
			closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeSaturday.Clock()
			openingTimeSaturday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
			closingTimeSaturday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
			if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
				closingTimeSaturday = closingTimeSaturday.AddDate(0, 0, 1)
			}
			if businessScheduleRes.SecondOpeningTimeSaturday == zeroTime {
				if startOrderTimeWeekday.Before(openingTimeSaturday) || endOrderTimeWeekday.After(closingTimeSaturday) {
					return errors.New("business closed")
				}
			} else {
				openingHour, openingMinutes, _ := businessScheduleRes.FirstOpeningTimeSaturday.Clock()
				closingHour, closingMinutes, _ := businessScheduleRes.FirstClosingTimeSaturday.Clock()
				secondOpeningTimeSaturday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), openingHour, openingMinutes, 0, 0, time.UTC)
				secondClosingTimeSaturday := time.Date(startOrderTimeWeekday.Year(), startOrderTimeWeekday.Month(), startOrderTimeWeekday.Day(), closingHour, closingMinutes, 0, 0, time.UTC)
				if openingHour > closingHour && closingHour >= 0 && closingHour <= 4 {
					secondClosingTimeSaturday = secondClosingTimeSaturday.AddDate(0, 0, 1)
				}
				if (startOrderTimeWeekday.Before(openingTimeSaturday) || (endOrderTimeWeekday.After(closingTimeSaturday) && startOrderTimeWeekday.Before(secondOpeningTimeSaturday))) || (startOrderTimeWeekday.After(secondClosingTimeSaturday) || endOrderTimeWeekday.After(secondClosingTimeSaturday)) {
					return errors.New("business closed")
				}
			}

		}
		_, err = i.dao.NewOrderedRepository().BatchCreateOrderedItem(tx, &orderedItems)
		if err != nil {
			return err
		}
		businessPaymentMethodId := uuid.MustParse(req.BusinessPaymentMethodId)
		businessPaymentMethod, err := i.dao.NewBusinessPaymentMethodRepository().GetBusinessPaymentMethod(ctx, tx, &entity.BusinessPaymentMethod{ID: &businessPaymentMethodId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("business payment method not found")
		} else if err != nil {
			return err
		}
		createOrderRes, err := i.dao.NewOrderRepository().CreateOrder(tx, &entity.Order{ItemsQuantity: quantity, BusinessThumbnail: businessRes.Thumbnail, OrderType: req.OrderType.String(), UserId: authorizationTokenRes.UserId, StartOrderTime: req.StartOrderTime.AsTime().UTC(), EndOrderTime: req.EndOrderTime.AsTime().UTC(), Coordinates: location, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: (*listCartItemRes)[0].BusinessId, PriceCup: price_cup.String(), CreateTime: createTime, UpdateTime: createTime, Number: userAddress.Number, Address: userAddress.Address, Instructions: req.Instructions, BusinessName: businessRes.Name, Status: "OrderStatusTypeOrdered", Phone: req.Phone, PaymentMethodType: businessPaymentMethod.Type, DeliveryPriceCup: businessRes.DeliveryPriceCup})
		if err != nil {
			return err
		}
		_, err = i.dao.NewOrderLifecycleRepository().CreateOrderLifecycle(tx, &entity.OrderLifecycle{Status: createOrderRes.Status, OrderId: createOrderRes.ID})
		if err != nil {
			return err
		}
		unionOrderAndOrderedItems := make([]entity.UnionOrderAndOrderedItem, 0, len(orderedItems))
		for _, item := range orderedItems {
			unionOrderAndOrderedItems = append(unionOrderAndOrderedItems, entity.UnionOrderAndOrderedItem{OrderId: createOrderRes.ID, OrderedItemId: item.ID})
		}
		_, err = i.dao.NewUnionOrderAndOrderedItemRepository().BatchCreateUnionOrderAndOrderedItem(tx, &unionOrderAndOrderedItems)
		if err != nil {
			return err
		}
		_, err = i.dao.NewCartItemRepository().DeleteCartItem(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, nil)
		if err != nil {
			return err
		}
		res = &pb.Order{BusinessName: businessRes.Name, BusinessThumbnail: i.config.BusinessAvatarBulkName + "/" + createOrderRes.BusinessThumbnail, ItemsQuantity: quantity, Status: *utils.ParseOrderStatusType(&createOrderRes.Status), OrderType: *utils.ParseOrderType(&createOrderRes.OrderType), Number: createOrderRes.Number, BusinessId: createOrderRes.BusinessId.String(), UserId: createOrderRes.UserId.String(), StartOrderTime: timestamppb.New(createOrderRes.StartOrderTime), EndOrderTime: timestamppb.New(createOrderRes.EndOrderTime), Coordinates: &pb.Point{Latitude: createOrderRes.Coordinates.FlatCoords()[0], Longitude: createOrderRes.Coordinates.FlatCoords()[1]}, PriceCup: price_cup.String(), CreateTime: timestamppb.New(createOrderRes.CreateTime), UpdateTime: timestamppb.New(createOrderRes.UpdateTime), Address: createOrderRes.Address, Instructions: createOrderRes.Instructions, Id: createOrderRes.ID.String(), ShortId: createOrderRes.ShortId, DeliveryPriceCup: createOrderRes.DeliveryPriceCup}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderService) ListOrder(ctx context.Context, req *pb.ListOrderRequest, md *utils.ClientMetadata) (*pb.ListOrderResponse, error) {
	var res pb.ListOrderResponse
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		var ordersRes *[]entity.OrderBusiness
		var nextPage time.Time
		if req.NextPage == nil {
			nextPage = time.Now()
		} else {
			nextPage = req.NextPage.AsTime()
		}
		// if !req.Upcoming && !req.History {
		// 	ordersRes, err = i.dao.NewOrderRepository().ListOrderWithBusiness(tx, &entity.OrderBusiness{CreateTime: nextPage, UserId: authorizationTokenRes.UserId})
		// } else {
		// 	ordersRes, err = i.dao.NewOrderRepository().ListOrderFilter(tx, &entity.OrderBusiness{CreateTime: nextPage, UserId: authorizationTokenRes.UserId}, req.Upcoming)
		// }
		ordersRes, err = i.dao.NewOrderRepository().ListOrderFilter(tx, &entity.OrderBusiness{CreateTime: nextPage, UserId: authorizationTokenRes.UserId}, req.Upcoming)
		if err != nil {
			return err
		}
		if len(*ordersRes) > 10 {
			*ordersRes = (*ordersRes)[:len(*ordersRes)-1]
			res.NextPage = timestamppb.New((*ordersRes)[len(*ordersRes)-1].CreateTime)
		} else if len(*ordersRes) > 0 {
			res.NextPage = timestamppb.New((*ordersRes)[len(*ordersRes)-1].CreateTime)
		} else {
			res.NextPage = timestamppb.New(nextPage)
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
				Instructions:      item.Instructions,
				UserId:            item.UserId.String(),
				StartOrderTime:    timestamppb.New(item.StartOrderTime),
				EndOrderTime:      timestamppb.New(item.EndOrderTime),
				Status:            *utils.ParseOrderStatusType(&item.Status),
				OrderType:         *utils.ParseOrderType(&item.OrderType),
				Coordinates:       &pb.Point{Latitude: item.Coordinates.Coords()[1], Longitude: item.Coordinates.Coords()[0]},
				BusinessId:        item.BusinessId.String(),
				BusinessThumbnail: i.config.BusinessAvatarBulkName + "/" + item.BusinessThumbnail,
				DeliveryPriceCup:  item.DeliveryPriceCup,
				CreateTime:        timestamppb.New(item.CreateTime),
				UpdateTime:        timestamppb.New(item.UpdateTime),
			})
		}
		res.Orders = ordersResponse
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
