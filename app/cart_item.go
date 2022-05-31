package app

import (
	"context"
	"strings"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *CartItemServer) EmptyCartItem(ctx context.Context, req *gp.Empty) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	err := m.cartItemService.EmptyCartItem(&dto.EmptyCartItemRequest{Metadata: &md})
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
}

func (m *CartItemServer) ListCartItem(ctx context.Context, req *pb.ListCartItemRequest) (*pb.ListCartItemResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.cartItemService.ListCartItem(ctx, req, md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *CartItemServer) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.AddCartItemResponse, error) {
	var invalidItemId, invalidLocation, invalidQuantity *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Quantity <= 0 {
		invalidArgs = true
		invalidQuantity = &epb.BadRequest_FieldViolation{
			Field:       "Quantity",
			Description: "The Quantity value must be greater than 0",
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
	if req.ItemId == "" {
		invalidArgs = true
		invalidItemId = &epb.BadRequest_FieldViolation{
			Field:       "ItemId",
			Description: "The ItemId field is required",
		}
	} else if req.ItemId != "" {
		if !utils.IsValidUUID(&req.ItemId) {
			invalidArgs = true
			invalidItemId = &epb.BadRequest_FieldViolation{
				Field:       "ItemId",
				Description: "The ItemId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		if invalidQuantity != nil {
			st, _ = st.WithDetails(
				invalidQuantity,
			)
		}
		if invalidItemId != nil {
			st, _ = st.WithDetails(
				invalidItemId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.cartItemService.AddCartItem(ctx, req, md)
	if err != nil {
		errorr := strings.Split(err.Error(), ":")
		switch errorr[0] {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "out of range":
			st = status.New(codes.InvalidArgument, "Out of range")
		case "no_availability":
			st = status.New(codes.InvalidArgument, "No availability")
			ds, _ := st.WithDetails(
				&epb.QuotaFailure{
					Violations: []*epb.QuotaFailure_Violation{{
						Subject:     "Availability",
						Description: errorr[2],
					}},
				},
			)
			st = ds
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *CartItemServer) ReduceCartItem(ctx context.Context, req *pb.ReduceCartItemRequest) (*pb.ReduceCartItemResponse, error) {
	var invalidItemId, invalidLocation, invalidQuantity *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Quantity <= 0 {
		invalidArgs = true
		invalidQuantity = &epb.BadRequest_FieldViolation{
			Field:       "Quantity",
			Description: "The Quantity value must be greater than 0",
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
	if req.ItemId == "" {
		invalidArgs = true
		invalidItemId = &epb.BadRequest_FieldViolation{
			Field:       "ItemId",
			Description: "The ItemId field is required",
		}
	} else if req.ItemId != "" {
		if !utils.IsValidUUID(&req.ItemId) {
			invalidArgs = true
			invalidItemId = &epb.BadRequest_FieldViolation{
				Field:       "ItemId",
				Description: "The ItemId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		if invalidQuantity != nil {
			st, _ = st.WithDetails(
				invalidQuantity,
			)
		}
		if invalidItemId != nil {
			st, _ = st.WithDetails(
				invalidItemId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.cartItemService.ReduceCartItem(ctx, req, md)
	if err != nil {
		errorr := strings.Split(err.Error(), ":")
		switch errorr[0] {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "cartitem not found":
			st = status.New(codes.NotFound, "CartItem not found")
		case "out of range":
			st = status.New(codes.InvalidArgument, "Out of range")
		case "no_availability":
			st = status.New(codes.InvalidArgument, "No availability")
			ds, _ := st.WithDetails(
				&epb.QuotaFailure{
					Violations: []*epb.QuotaFailure_Violation{{
						Subject:     "Availability",
						Description: errorr[2],
					}},
				},
			)
			st = ds
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *CartItemServer) DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest) (*gp.Empty, error) {
	var invalidId, invalidLocation *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
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
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.cartItemService.DeleteCartItem(ctx, req, md)
	if err != nil {
		errorr := strings.Split(err.Error(), ":")
		switch errorr[0] {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "cartitem not found":
			st = status.New(codes.NotFound, "CartItem not found")
		case "out of range":
			st = status.New(codes.InvalidArgument, "Out of range")
		case "no_availability":
			st = status.New(codes.InvalidArgument, "No availability")
			ds, _ := st.WithDetails(
				&epb.QuotaFailure{
					Violations: []*epb.QuotaFailure_Violation{{
						Subject:     "Availability",
						Description: errorr[2],
					}},
				},
			)
			st = ds
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *CartItemServer) CartItemIsEmpty(ctx context.Context, req *gp.Empty) (*pb.CartItemIsEmptyResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.cartItemService.CartItemIsEmpty(ctx, req, md)
	if err != nil {
		errorr := strings.Split(err.Error(), ":")
		switch errorr[0] {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "cartitem not found":
			st = status.New(codes.NotFound, "CartItem not found")
		case "out of range":
			st = status.New(codes.InvalidArgument, "Out of range")
		case "no_availability":
			st = status.New(codes.InvalidArgument, "No availability")
			ds, _ := st.WithDetails(
				&epb.QuotaFailure{
					Violations: []*epb.QuotaFailure_Violation{{
						Subject:     "Availability",
						Description: errorr[2],
					}},
				},
			)
			st = ds
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
