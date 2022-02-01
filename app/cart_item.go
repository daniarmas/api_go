package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
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
