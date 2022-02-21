package app

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *OrderServer) ListOrder(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	var nextPage time.Time
	if req.NextPage.Nanos == 0 && req.NextPage.Seconds == 0 {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	listOrderResponse, err := m.orderService.ListOrder(&dto.ListOrderRequest{NextPage: nextPage, Metadata: &md})
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
	ordersResponse := make([]*pb.Order, 0, len(*listOrderResponse.Orders))
	for _, item := range *listOrderResponse.Orders {
		ordersResponse = append(ordersResponse, &pb.Order{
			Id:             item.ID.String(),
			Price:          item.Price,
			BuildingNumber: item.BuildingNumber,
			HouseNumber:    item.HouseNumber,
			UserFk:         item.UserFk.String(),
			AppVersion:     item.AppVersion,
			DeviceFk:       item.DeviceFk.String(),
			DeliveryDate:   timestamppb.New(item.DeliveryDate),
			Status:         *utils.ParseOrderStatusType(&item.Status),
			DeliveryType:   *utils.ParseDeliveryType(&item.DeliveryType),
			ResidenceType:  *utils.ParseOrderResidenceType(&item.ResidenceType),
			Coordinates:    &pb.Point{Latitude: item.Coordinates.Coords()[1], Longitude: item.Coordinates.Coords()[0]},
			BusinessFk:     item.BusinessFk.String(),
			CreateTime:     timestamppb.New(item.CreateTime),
			UpdateTime:     timestamppb.New(item.UpdateTime),
		})
	}
	return &pb.ListOrderResponse{Orders: ordersResponse, NextPage: timestamppb.New(listOrderResponse.NextPage)}, nil
}