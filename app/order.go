package app

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *OrderServer) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	var nextPage time.Time
	if req.NextPage.Nanos == 0 && req.NextPage.Seconds == 0 {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	listOrderResponse, err := m.orderService.ListOrder(&dto.ListOrderRequest{NextPage: nextPage, Metadata: &md})
	if err != nil {
		switch err.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "business is open":
			st = status.New(codes.InvalidArgument, "Business is open")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		case "item in the cart":
			st = status.New(codes.InvalidArgument, "Item in the cart")
		case "cartitem not found":
			st = status.New(codes.NotFound, "CartItem not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	ordersResponse := make([]*pb.Order, 0, len(*listOrderResponse.Orders))
	for _, item := range *listOrderResponse.Orders {
		ordersResponse = append(ordersResponse, &pb.Order{
			Id:             item.ID.String(),
			BusinessName:   item.BusinessName,
			Quantity:       item.Quantity,
			Price:          item.Price,
			BuildingNumber: item.BuildingNumber,
			HouseNumber:    item.HouseNumber,
			UserFk:         item.UserFk.String(),
			OrderDate:      timestamppb.New(item.OrderDate),
			Status:         *utils.ParseOrderStatusType(&item.Status),
			OrderType:      *utils.ParseDeliveryType(&item.OrderType),
			ResidenceType:  *utils.ParseOrderResidenceType(&item.ResidenceType),
			Coordinates:    &pb.Point{Latitude: item.Coordinates.Coords()[1], Longitude: item.Coordinates.Coords()[0]},
			BusinessFk:     item.BusinessFk.String(),
			CreateTime:     timestamppb.New(item.CreateTime),
			UpdateTime:     timestamppb.New(item.UpdateTime),
		})
	}
	return &pb.ListOrderResponse{Orders: ordersResponse, NextPage: timestamppb.New(listOrderResponse.NextPage)}, nil
}

func (m *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	cartItems := make([]uuid.UUID, 0, len(req.CartItems))
	for _, item := range req.CartItems {
		cartItems = append(cartItems, uuid.MustParse(item))
	}
	createOrderRes, createOrderErr := m.orderService.CreateOrder(&dto.CreateOrderRequest{CartItems: &cartItems, Status: req.Status.String(), OrderType: req.OrderType.String(), ResidenceType: req.ResidenceType.String(), BuildingNumber: req.BuildingNumber, HouseNumber: req.HouseNumber, OrderDate: req.OrderDate.AsTime(), Coordinates: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Coordinates.Latitude, req.Coordinates.Longitude}).SetSRID(4326)}, Metadata: &md})
	if createOrderErr != nil {
		switch createOrderErr.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "business is open":
			st = status.New(codes.InvalidArgument, "Business is open")
		case "item in the cart":
			st = status.New(codes.InvalidArgument, "Item in the cart")
		case "cartitem not found":
			st = status.New(codes.NotFound, "CartItem not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.CreateOrderResponse{Order: &pb.Order{Id: createOrderRes.Order.ID.String(), BuildingNumber: createOrderRes.Order.BuildingNumber, Price: createOrderRes.Order.Price, UserFk: createOrderRes.Order.UserFk.String(), BusinessFk: createOrderRes.Order.BusinessFk.String(), Status: *utils.ParseOrderStatusType(&createOrderRes.Order.Status), OrderType: *utils.ParseDeliveryType(&createOrderRes.Order.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&createOrderRes.Order.ResidenceType), HouseNumber: createOrderRes.Order.HouseNumber, CreateTime: timestamppb.New(createOrderRes.Order.CreateTime), UpdateTime: timestamppb.New(createOrderRes.Order.UpdateTime), OrderDate: timestamppb.New(createOrderRes.Order.OrderDate), Coordinates: &pb.Point{Latitude: createOrderRes.Order.Coordinates.FlatCoords()[0], Longitude: createOrderRes.Order.Coordinates.FlatCoords()[1]}}}, nil
}

func (m *OrderServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	updateOrderRes, updateOrderErr := m.orderService.UpdateOrder(&dto.UpdateOrderRequest{Id: uuid.MustParse(req.Id), Status: req.Status.String(), Metadata: &md})
	if updateOrderErr != nil {
		switch updateOrderErr.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "invalid status value":
			st = status.New(codes.InvalidArgument, "Invalid status value")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.UpdateOrderResponse{Order: &pb.Order{Id: updateOrderRes.Order.ID.String(), BuildingNumber: updateOrderRes.Order.BuildingNumber, Price: updateOrderRes.Order.Price, UserFk: updateOrderRes.Order.UserFk.String(), BusinessFk: updateOrderRes.Order.BusinessFk.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Order.Status), OrderType: *utils.ParseDeliveryType(&updateOrderRes.Order.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.Order.ResidenceType), HouseNumber: updateOrderRes.Order.HouseNumber, CreateTime: timestamppb.New(updateOrderRes.Order.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.Order.UpdateTime), OrderDate: timestamppb.New(updateOrderRes.Order.OrderDate), Coordinates: &pb.Point{Latitude: updateOrderRes.Order.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Order.Coordinates.FlatCoords()[1]}}}, nil
}

func (m *OrderServer) ListOrderedItem(ctx context.Context, req *pb.ListOrderedItemRequest) (*pb.ListOrderedItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	listOrderedItemRes, listOrderedItemErr := m.orderService.ListOrderedItemWithItem(&dto.ListOrderedItemRequest{OrderFk: uuid.MustParse(req.OrderFk), Metadata: &md})
	if listOrderedItemErr != nil {
		switch listOrderedItemErr.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "invalid status value":
			st = status.New(codes.InvalidArgument, "Invalid status value")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	orderedItems := make([]*pb.OrderedItem, 0, len(*listOrderedItemRes.OrderedItems))
	for _, item := range *listOrderedItemRes.OrderedItems {
		orderedItems = append(orderedItems, &pb.OrderedItem{Id: item.ID.String(), Name: item.Name, Price: item.Price, ItemFk: item.ItemFk.String(), Quantity: item.Quantity, UserFk: item.UserFk.String(), CreateTime: timestamppb.New(item.CreateTime), UpdateTime: timestamppb.New(item.UpdateTime), CartItemFk: item.CartItemFk.String()})
	}
	return &pb.ListOrderedItemResponse{OrderedItems: orderedItems}, nil
}
