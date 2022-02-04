package app

import (
	"context"
	"strings"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (m *CartItemServer) ListCartItem(ctx context.Context, req *pb.ListCartItemRequest) (*pb.ListCartItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	listCartItemsResponse, err := m.cartItemService.ListCartItemAndItem(&dto.ListCartItemRequest{NextPage: req.NextPage, Metadata: &md})
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
			CreateTime:           item.CreateTime.String(),
			UpdateTime:           item.UpdateTime.String(),
			Cursor:               item.Cursor,
			ItemFk:               item.ItemFk.String(),
			AuthorizationTokenFk: item.AuthorizationTokenFk.String(),
			Quantity:             item.Quantity,
			Thumbnail:            item.Thumbnail,
			ThumbnailBlurHash:    item.ThumbnailBlurHash,
		})
	}
	return &pb.ListCartItemResponse{CartItems: itemsResponse, NextPage: listCartItemsResponse.NextPage}, nil
}

func (m *CartItemServer) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.AddCartItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	cartItemsResponse, err := m.cartItemService.AddCartItem(&dto.AddCartItem{ItemFk: req.ItemFk, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, Metadata: &md, Quantity: req.Quantity})
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
		CreateTime:           cartItemsResponse.CreateTime.String(),
		UpdateTime:           cartItemsResponse.UpdateTime.String(),
		Cursor:               cartItemsResponse.Cursor,
		ItemFk:               cartItemsResponse.ItemFk.String(),
		AuthorizationTokenFk: cartItemsResponse.AuthorizationTokenFk.String(),
		Quantity:             cartItemsResponse.Quantity,
		// Thumbnail:            cartItemsResponse.Thumbnail,
		// ThumbnailBlurHash:    cartItemsResponse.ThumbnailBlurHash,
	}}, nil
}
