package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *OrderServer) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.orderService.ListOrder(ctx, req, md)
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
	return res, nil
}

func (m *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var invalidCartItems, invalidLocation, invalidOrderType, invalidResidenceType, invalidNumber, invalidAddress *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if len(req.CartItems) == 0 {
		invalidArgs = true
		invalidCartItems = &epb.BadRequest_FieldViolation{
			Field:       "CartItems",
			Description: "The CartItems field is required",
		}
	} else {
		for _, elem := range req.CartItems {
			if !utils.IsValidUUID(&elem) {
				invalidArgs = true
				invalidCartItems = &epb.BadRequest_FieldViolation{
					Field:       "CartItems",
					Description: "The CartItems contains not a valid uuid v4",
				}
				break
			}
		}
	}
	if req.Location == nil {
		invalidArgs = true
		invalidLocation = &epb.BadRequest_FieldViolation{
			Field:       "Location",
			Description: "The Location field is required",
		}
	} else if req.Location != nil {
		if req.Location.Latitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "Location.Latitude",
				Description: "The Location.Latitude field is required",
			}
		} else if req.Location.Longitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "Location.Longitude",
				Description: "The Location.Longitude field is required",
			}
		}
	}
	if req.Address == "" {
		invalidArgs = true
		invalidAddress = &epb.BadRequest_FieldViolation{
			Field:       "Address",
			Description: "The Address field is required",
		}
	}
	if req.Number == "" {
		invalidArgs = true
		invalidNumber = &epb.BadRequest_FieldViolation{
			Field:       "Number",
			Description: "The Number field is required",
		}
	}
	if req.OrderType == *pb.OrderType_OrderTypeUnspecified.Enum() {
		invalidArgs = true
		invalidOrderType = &epb.BadRequest_FieldViolation{
			Field:       "OrderType",
			Description: "The OrderType field is required",
		}
	}
	if req.ResidenceType == *pb.ResidenceType_ResidenceTypeUnspecified.Enum() {
		invalidArgs = true
		invalidResidenceType = &epb.BadRequest_FieldViolation{
			Field:       "ResidenceType",
			Description: "The ResidenceType field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		if invalidAddress != nil {
			st, _ = st.WithDetails(
				invalidAddress,
			)
		}
		if invalidNumber != nil {
			st, _ = st.WithDetails(
				invalidNumber,
			)
		}
		if invalidCartItems != nil {
			st, _ = st.WithDetails(
				invalidCartItems,
			)
		}
		if invalidResidenceType != nil {
			st, _ = st.WithDetails(
				invalidResidenceType,
			)
		}
		if invalidOrderType != nil {
			st, _ = st.WithDetails(
				invalidOrderType,
			)
		}
		return nil, st.Err()
	}
	res, createOrderErr := m.orderService.CreateOrder(ctx, req, md)
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
		case "business closed":
			st = status.New(codes.InvalidArgument, "Business closed")
		case "invalid schedule":
			st = status.New(codes.InvalidArgument, "Invalid schedule")
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
	return res, nil
}

func (m *OrderServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	id := uuid.MustParse(req.Id)
	updateOrderRes, updateOrderErr := m.orderService.UpdateOrder(&dto.UpdateOrderRequest{Id: &id, Status: req.Status.String(), Metadata: &md})
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
	return &pb.UpdateOrderResponse{Order: &pb.Order{Id: updateOrderRes.Order.ID.String(), Address: updateOrderRes.Order.Address, Instructions: updateOrderRes.Order.Instructions, Price: updateOrderRes.Order.Price, UserId: updateOrderRes.Order.UserId.String(), BusinessId: updateOrderRes.Order.BusinessId.String(), Status: *utils.ParseOrderStatusType(&updateOrderRes.Order.Status), OrderType: *utils.ParseOrderType(&updateOrderRes.Order.OrderType), ResidenceType: *utils.ParseOrderResidenceType(&updateOrderRes.Order.ResidenceType), Number: updateOrderRes.Order.Number, CreateTime: timestamppb.New(updateOrderRes.Order.CreateTime), UpdateTime: timestamppb.New(updateOrderRes.Order.UpdateTime), OrderDate: timestamppb.New(updateOrderRes.Order.OrderDate), Coordinates: &pb.Point{Latitude: updateOrderRes.Order.Coordinates.FlatCoords()[0], Longitude: updateOrderRes.Order.Coordinates.FlatCoords()[1]}}}, nil
}

func (m *OrderServer) ListOrderedItem(ctx context.Context, req *pb.ListOrderedItemRequest) (*pb.ListOrderedItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	orderId := uuid.MustParse(req.OrderId)
	listOrderedItemRes, listOrderedItemErr := m.orderService.ListOrderedItemWithItem(&dto.ListOrderedItemRequest{OrderId: &orderId, Metadata: &md})
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
		orderedItems = append(orderedItems, &pb.OrderedItem{Id: item.ID.String(), Name: item.Name, Price: item.Price, ItemId: item.ItemId.String(), Quantity: item.Quantity, UserId: item.UserId.String(), CreateTime: timestamppb.New(item.CreateTime), UpdateTime: timestamppb.New(item.UpdateTime), CartItemId: item.CartItemId.String()})
	}
	return &pb.ListOrderedItemResponse{OrderedItems: orderedItems}, nil
}
