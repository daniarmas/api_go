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

func (m *PaymentMethodServer) ListBusinessPaymentMethod(ctx context.Context, req *pb.ListBusinessPaymentMethodRequest) (*pb.ListBusinessPaymentMethodResponse, error) {
	var invalidBusinessId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.BusinessId == "" {
		invalidArgs = true
		invalidBusinessId = &epb.BadRequest_FieldViolation{
			Field:       "businessId",
			Description: "The businessId field is required",
		}
	} else if req.BusinessId != "" {
		if !utils.IsValidUUID(&req.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "businessId",
				Description: "The businessId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidBusinessId != nil {
			st, _ = st.WithDetails(
				invalidBusinessId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.paymentMethodService.ListBusinessPaymentMethod(ctx, req, md)
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

func (m *PaymentMethodServer) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest) (*gp.Empty, error) {
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
			Description: "The id field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "Id",
				Description: "The id field is not a valid uuid v4",
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
	res, err := m.paymentMethodService.DeletePaymentMethod(ctx, req, md)
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "payment method not found":
			st = status.New(codes.NotFound, "Payment Method not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *PaymentMethodServer) ListPaymentMethod(ctx context.Context, req *pb.ListPaymentMethodRequest) (*pb.ListPaymentMethodResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	res, err := m.paymentMethodService.ListPaymentMethod(ctx, req, md)
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

func (m *PaymentMethodServer) UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest) (*pb.PaymentMethod, error) {
	var invalidId, invalidType *epb.BadRequest_FieldViolation
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
			Description: "The id field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "Id",
				Description: "The id field is not a valid uuid v4",
			}
		}
	}
	if req.PaymentMethod.Type == pb.PaymentMethodType_PaymentMethodTypeUnspecified {
		invalidArgs = true
		invalidType = &epb.BadRequest_FieldViolation{
			Field:       "paymentMethod.type",
			Description: "The paymentMethod.type field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		if invalidType != nil {
			st, _ = st.WithDetails(
				invalidType,
			)
		}
		return nil, st.Err()
	}
	res, err := m.paymentMethodService.UpdatePaymentMethod(ctx, req, md)
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "payment method not found":
			st = status.New(codes.NotFound, "Payment Method not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *PaymentMethodServer) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest) (*pb.PaymentMethod, error) {
	var invalidType *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.PaymentMethod.Type == pb.PaymentMethodType_PaymentMethodTypeUnspecified {
		invalidArgs = true
		invalidType = &epb.BadRequest_FieldViolation{
			Field:       "paymentMethod.type",
			Description: "The paymentMethod.type field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidType != nil {
			st, _ = st.WithDetails(
				invalidType,
			)
		}
		return nil, st.Err()
	}
	res, err := m.paymentMethodService.CreatePaymentMethod(ctx, req, md)
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
