package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		case "cart items not found":
			st = status.New(codes.InvalidArgument, "Cart items not found")
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
	var invalidId, invalidStatus *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Id == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "Id",
			Description: "The Id field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "Id",
				Description: "The Id field is not a valid uuid v4",
			}
		}
	}
	if req.Status == *pb.OrderStatusType_OrderStatusTypeUnspecified.Enum() {
		invalidArgs = true
		invalidStatus = &epb.BadRequest_FieldViolation{
			Field:       "Status",
			Description: "The Status field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		if invalidStatus != nil {
			st, _ = st.WithDetails(
				invalidStatus,
			)
		}
		return nil, st.Err()
	}
	res, updateOrderErr := m.orderService.UpdateOrder(ctx, req, md)
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
	return res, nil
}

func (m *OrderServer) ListOrderedItem(ctx context.Context, req *pb.ListOrderedItemRequest) (*pb.ListOrderedItemResponse, error) {
	var invalidOrderId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.OrderId == "" {
		invalidArgs = true
		invalidOrderId = &epb.BadRequest_FieldViolation{
			Field:       "OrderId",
			Description: "The OrderId field is required",
		}
	} else if req.OrderId != "" {
		if !utils.IsValidUUID(&req.OrderId) {
			invalidArgs = true
			invalidOrderId = &epb.BadRequest_FieldViolation{
				Field:       "OrderId",
				Description: "The OrderId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidOrderId != nil {
			st, _ = st.WithDetails(
				invalidOrderId,
			)
		}
		return nil, st.Err()
	}
	res, listOrderedItemErr := m.orderService.ListOrderedItemWithItem(ctx, req, md)
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
	return res, nil
}
