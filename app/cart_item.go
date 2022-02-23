package app

import (
	"context"
	"strings"
	"time"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
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

func (m *CartItemServer) ListCartItem(ctx context.Context, req *pb.ListCartItemRequest) (*pb.ListCartItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	var nextPage = req.NextPage.AsTime()
	if req.NextPage.Nanos == 0 && req.NextPage.Seconds == 0 {
		nextPage = time.Now()
	}
	listCartItemsResponse, err := m.cartItemService.ListCartItemAndItem(&dto.ListCartItemRequest{NextPage: nextPage, Metadata: &md, MunicipalityFk: uuid.MustParse(req.MunicipalityFk)})
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
	itemsResponse := make([]*pb.CartItem, 0, len(listCartItemsResponse.CartItems))
	for _, item := range listCartItemsResponse.CartItems {
		itemsResponse = append(itemsResponse, &pb.CartItem{
			Id:                   item.ID.String(),
			Name:                 item.Name,
			Price:                item.Price,
			ItemFk:               item.ItemFk.String(),
			AuthorizationTokenFk: item.AuthorizationTokenFk.String(),
			Quantity:             item.Quantity,
			Thumbnail:            item.Thumbnail,
			ThumbnailBlurHash:    item.ThumbnailBlurHash,
			CreateTime:           timestamppb.New(item.CreateTime),
			UpdateTime:           timestamppb.New(item.UpdateTime),
		})
	}
	return &pb.ListCartItemResponse{CartItems: itemsResponse, NextPage: timestamppb.New(listCartItemsResponse.NextPage)}, nil
}

func (m *CartItemServer) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.AddCartItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	cartItemsResponse, err := m.cartItemService.AddCartItem(&dto.AddCartItem{ItemFk: req.ItemFk, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, Quantity: req.Quantity, MunicipalityFk: uuid.MustParse(req.MunicipalityFk)})
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
		ItemFk:               cartItemsResponse.ItemFk.String(),
		AuthorizationTokenFk: cartItemsResponse.AuthorizationTokenFk.String(),
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
	cartItemsResponse, err := m.cartItemService.ReduceCartItem(&dto.ReduceCartItem{ItemFk: req.ItemFk, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, MunicipalityFk: uuid.MustParse(req.MunicipalityFk)})
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
			ItemFk:               cartItemsResponse.ItemFk.String(),
			AuthorizationTokenFk: cartItemsResponse.AuthorizationTokenFk.String(),
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
	err := m.cartItemService.DeleteCartItem(&dto.DeleteCartItemRequest{CartItemFk: req.Id, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, MunicipalityFk: uuid.MustParse(req.MunicipalityFk)})
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
