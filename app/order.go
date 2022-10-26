package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *OrderServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*gp.Empty, error) {
	var invalidId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
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
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.orderService.CancelOrder(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "order not found":
			st = status.New(codes.NotFound, "Order not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *OrderServer) GetCheckoutInfo(ctx context.Context, req *pb.GetCheckoutInfoRequest) (*pb.GetCheckoutInfoResponse, error) {
	var invalidId, invalidLocation *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.BusinessId == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "businessId",
			Description: "The businessId field is required",
		}
	} else if req.BusinessId != "" {
		if !utils.IsValidUUID(&req.BusinessId) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "businessId",
				Description: "The businessId field is not a valid uuid v4",
			}
		}
	}
	if req.Coordinates == nil {
		invalidArgs = true
		invalidLocation = &epb.BadRequest_FieldViolation{
			Field:       "coordinates",
			Description: "The coordinates field is required",
		}
	} else if req.Coordinates != nil {
		if req.Coordinates.Latitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "coordinates.latitude",
				Description: "The coordinates.latitude field is required",
			}
		} else if req.Coordinates.Longitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "coordinates.longitude",
				Description: "The coordinates.longitude field is required",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		return nil, st.Err()
	}
	res, err := m.orderService.GetCheckoutInfo(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "business not in range":
			st = status.New(codes.InvalidArgument, "Business not in range")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "business not found":
			st = status.New(codes.NotFound, "Business not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	var invalidId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
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
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.orderService.GetOrder(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "order not found":
			st = status.New(codes.NotFound, "Order not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *OrderServer) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	res, err := m.orderService.ListOrder(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	var invalidCartItems, invalidUserAddressId, invalidPhone, invalidBusinessPaymentMethodId, invalidOrderType *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.UserAddressId == "" && req.OrderType == pb.OrderType_OrderTypeHomeDelivery {
		invalidArgs = true
		invalidUserAddressId = &epb.BadRequest_FieldViolation{
			Field:       "userAddressId",
			Description: "The userAddressId field is required",
		}
	} else if req.UserAddressId != "" && req.OrderType == pb.OrderType_OrderTypeHomeDelivery {
		if !utils.IsValidUUID(&req.UserAddressId) {
			invalidArgs = true
			invalidUserAddressId = &epb.BadRequest_FieldViolation{
				Field:       "userAddressId",
				Description: "The userAddressId field is not a valid uuid v4",
			}
		}
	}
	if req.BusinessPaymentMethodId == "" {
		invalidArgs = true
		invalidBusinessPaymentMethodId = &epb.BadRequest_FieldViolation{
			Field:       "businessPaymentMethodId",
			Description: "The businessPaymentMethodId field is required",
		}
	} else if req.BusinessPaymentMethodId != "" {
		if !utils.IsValidUUID(&req.BusinessPaymentMethodId) {
			invalidArgs = true
			invalidBusinessPaymentMethodId = &epb.BadRequest_FieldViolation{
				Field:       "businessPaymentMethodId",
				Description: "The businessPaymentMethodId field is not a valid uuid v4",
			}
		}
	}
	if req.Phone == "" {
		invalidArgs = true
		invalidPhone = &epb.BadRequest_FieldViolation{
			Field:       "phone",
			Description: "The phone field is required",
		}
	}
	if req.OrderType == *pb.OrderType_OrderTypeUnspecified.Enum() {
		invalidArgs = true
		invalidOrderType = &epb.BadRequest_FieldViolation{
			Field:       "OrderType",
			Description: "The OrderType field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidBusinessPaymentMethodId != nil {
			st, _ = st.WithDetails(
				invalidBusinessPaymentMethodId,
			)
		}
		if invalidPhone != nil {
			st, _ = st.WithDetails(
				invalidPhone,
			)
		}
		if invalidUserAddressId != nil {
			st, _ = st.WithDetails(
				invalidUserAddressId,
			)
		}
		if invalidCartItems != nil {
			st, _ = st.WithDetails(
				invalidCartItems,
			)
		}
		if invalidOrderType != nil {
			st, _ = st.WithDetails(
				invalidOrderType,
			)
		}
		return nil, st.Err()
	}
	res, err := m.orderService.CreateOrder(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "not fulfilled the previous time of the business":
			st = status.New(codes.InvalidArgument, "Not fulfilled the previous time of the business")
		case "business closed":
			st = status.New(codes.InvalidArgument, "Business closed")
		case "cart items not found":
			st = status.New(codes.InvalidArgument, "Cart items not found")
		case "user address not found":
			st = status.New(codes.NotFound, "User address not found")
		case "business payment method not found":
			st = status.New(codes.NotFound, "Business payment method not found")
		case "business not in range":
			st = status.New(codes.InvalidArgument, "Business not in range")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *OrderServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.Order, error) {
	var invalidId, invalidStatus *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.Order.Id == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "Id",
			Description: "The Id field is required",
		}
	} else if req.Order.Id != "" {
		if !utils.IsValidUUID(&req.Order.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "Id",
				Description: "The Id field is not a valid uuid v4",
			}
		}
	}
	if req.Order.Status == *pb.OrderStatusType_OrderStatusTypeUnspecified.Enum() {
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
	res, err := m.orderService.UpdateOrder(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "order not found":
			st = status.New(codes.NotFound, "Order not found")
		case "status error":
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
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
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
	res, err := m.orderService.ListOrderedItemWithItem(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
