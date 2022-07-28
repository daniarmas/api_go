package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *PaymentMethodServer) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest) (*pb.PaymentMethod, error) {
	var invalidName, invalidType *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.PaymentMethod.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "paymentMethod.name",
			Description: "The paymentMethod.name field is required",
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
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
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
		case "user configuration not found":
			st = status.New(codes.Unauthenticated, "User address not found")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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
