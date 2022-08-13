package app

import (
	"context"
	"strings"

	pb "github.com/daniarmas/api_go/pkg/grpc"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *CartItemServer) EmptyAndAddCartItem(ctx context.Context, req *pb.EmptyAndAddCartItemRequest) (*pb.CartItem, error) {
	var invalidItemId, invalidLocation, invalidQuantity *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
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
	res, err := m.cartItemService.EmptyAndAddCartItem(ctx, req, md)
	if err != nil {
		errorr := strings.Split(err.Error(), ":")
		switch errorr[0] {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "item not found":
			st = status.New(codes.NotFound, "Item not found")
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

func (m *CartItemServer) EmptyCartItem(ctx context.Context, req *gp.Empty) (*gp.Empty, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.cartItemService.EmptyCartItem(ctx, md)
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
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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

func (m *CartItemServer) ListCartItem(ctx context.Context, req *pb.ListCartItemRequest) (*pb.ListCartItemResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.cartItemService.ListCartItem(ctx, req, md)
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
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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

func (m *CartItemServer) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.CartItem, error) {
	var invalidItemId, invalidLocation, invalidQuantity *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "item not found":
			st = status.New(codes.NotFound, "Item not found")
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
		case "the items in the cart can only be from one business":
			st = status.New(codes.Unauthenticated, "The items in the cart can only be from one business")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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

func (m *CartItemServer) DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest) (*gp.Empty, error) {
	var invalidId, invalidItemId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.Id == "" && req.ItemId == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "Id|ItemId",
			Description: "The Id or ItemId fields are required",
		}
	}
	if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "Id",
				Description: "The Id field is not a valid uuid v4",
			}
		}
	}
	if req.ItemId != "" {
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
		if invalidItemId != nil {
			st, _ = st.WithDetails(
				invalidItemId,
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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

func (m *CartItemServer) IsEmptyCartItem(ctx context.Context, req *gp.Empty) (*pb.IsEmptyCartItemResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.cartItemService.IsEmptyCartItem(ctx, req, md)
	if err != nil {
		errorr := strings.Split(err.Error(), ":")
		switch errorr[0] {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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
