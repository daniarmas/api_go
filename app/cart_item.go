package app

import (
	"context"
	"strings"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	utils "github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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
	var invalidMunicipalityId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.MunicipalityId == "" {
		invalidArgs = true
		invalidMunicipalityId = &epb.BadRequest_FieldViolation{
			Field:       "MunicipalityId",
			Description: "The MunicipalityId field is required",
		}
	} else if req.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "MunicipalityId",
				Description: "The MunicipalityId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidMunicipalityId != nil {
			st, _ = st.WithDetails(
				invalidMunicipalityId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.cartItemService.ListCartItemAndItem(ctx, req, md)
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
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	municipalityId := uuid.MustParse(req.MunicipalityId)
	cartItemsResponse, err := m.cartItemService.AddCartItem(&dto.AddCartItem{ItemId: req.ItemId, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, Quantity: req.Quantity, MunicipalityId: &municipalityId})
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
	return &pb.AddCartItemResponse{CartItem: &pb.CartItem{
		Id:                   cartItemsResponse.ID.String(),
		Name:                 cartItemsResponse.Name,
		Price:                cartItemsResponse.Price,
		ItemId:               cartItemsResponse.ItemId.String(),
		AuthorizationTokenId: cartItemsResponse.AuthorizationTokenId.String(),
		Quantity:             cartItemsResponse.Quantity,
		CreateTime:           timestamppb.New(cartItemsResponse.CreateTime),
		UpdateTime:           timestamppb.New(cartItemsResponse.UpdateTime),
		// Thumbnail:            cartItemsResponse.Thumbnail,
		// ThumbnailBlurHash:    cartItemsResponse.ThumbnailBlurHash,
	}}, nil
}

func (m *CartItemServer) ReduceCartItem(ctx context.Context, req *pb.ReduceCartItemRequest) (*pb.ReduceCartItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	municipalityId := uuid.MustParse(req.MunicipalityId)
	cartItemsResponse, err := m.cartItemService.ReduceCartItem(&dto.ReduceCartItem{ItemId: req.ItemId, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, MunicipalityId: &municipalityId})
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
	if cartItemsResponse != nil {
		return &pb.ReduceCartItemResponse{CartItem: &pb.CartItem{
			Id:                   cartItemsResponse.ID.String(),
			Name:                 cartItemsResponse.Name,
			Price:                cartItemsResponse.Price,
			ItemId:               cartItemsResponse.ItemId.String(),
			AuthorizationTokenId: cartItemsResponse.AuthorizationTokenId.String(),
			Quantity:             cartItemsResponse.Quantity,
			CreateTime:           timestamppb.New(cartItemsResponse.CreateTime),
			UpdateTime:           timestamppb.New(cartItemsResponse.UpdateTime),
			// Thumbnail:            cartItemsResponse.Thumbnail,
			// ThumbnailBlurHash:    cartItemsResponse.ThumbnailBlurHash,
		}}, nil
	} else {
		return &pb.ReduceCartItemResponse{CartItem: &pb.CartItem{}}, nil
	}
}

func (m *CartItemServer) DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	municipalityId := uuid.MustParse(req.MunicipalityId)
	err := m.cartItemService.DeleteCartItem(&dto.DeleteCartItemRequest{CartItemId: req.Id, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, MunicipalityId: &municipalityId})
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
	return &gp.Empty{}, nil
}

func (m *CartItemServer) CartItemQuantity(ctx context.Context, req *gp.Empty) (*pb.CartItemQuantityResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	res, err := m.cartItemService.CartItemQuantity(&dto.CartItemQuantity{Metadata: &md})
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
	return &pb.CartItemQuantityResponse{IsFull: *res}, nil
}
