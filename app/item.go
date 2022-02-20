package app

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	ut "github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *ItemServer) ListItem(ctx context.Context, req *pb.ListItemRequest) (*pb.ListItemResponse, error) {
	var nextPage time.Time
	if req.NextPage.Nanos == 0 && req.NextPage.Seconds == 0 {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	listItemsResponse, err := m.itemService.ListItem(&dto.ListItemRequest{BusinessFk: req.BusinessFk, BusinessItemCategoryFk: req.ItemCategoryFk, NextPage: nextPage})
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*pb.Item, 0, len(listItemsResponse.Items))
	for _, item := range listItemsResponse.Items {
		itemsResponse = append(itemsResponse, &pb.Item{
			Id:                       item.ID.String(),
			Name:                     item.Name,
			Description:              item.Description,
			Price:                    item.Price,
			Availability:             int32(item.Availability),
			BusinessFk:               item.BusinessFk.String(),
			BusinessItemCategoryFk:   item.BusinessItemCategoryFk.String(),
			HighQualityPhoto:         item.HighQualityPhoto,
			HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
			LowQualityPhoto:          item.LowQualityPhoto,
			LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
			Thumbnail:                item.Thumbnail,
			ThumbnailBlurHash:        item.ThumbnailBlurHash,
			Cursor:                   int32(item.Cursor),
			CreateTime:               timestamppb.New(item.CreateTime),
			UpdateTime:               timestamppb.New(item.UpdateTime),
		})
	}
	return &pb.ListItemResponse{Items: itemsResponse, NextPage: timestamppb.New(listItemsResponse.NextPage)}, nil
}

func (m *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	var st *status.Status
	item, err := m.itemService.GetItem(&dto.GetItemRequest{Id: req.Id, Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}})
	if err != nil {
		switch err.Error() {
		case "record not found":
			st = status.New(codes.NotFound, "Item not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	// itemPhotos := make([]*pb.ItemPhoto, 0, len(item.ItemPhoto))
	// for _, e := range item.ItemPhoto {
	// 	itemPhotos = append(itemPhotos, &pb.ItemPhoto{
	// 		Id:                       e.ID.String(),
	// 		ItemFk:                   e.ItemFk.String(),
	// 		HighQualityPhoto:         e.HighQualityPhoto,
	// 		HighQualityPhotoBlurHash: e.HighQualityPhotoBlurHash,
	// 		LowQualityPhoto:          e.LowQualityPhoto,
	// 		LowQualityPhotoBlurHash:  e.LowQualityPhotoBlurHash,
	// 		Thumbnail:                e.Thumbnail,
	// 		ThumbnailBlurHash:        e.ThumbnailBlurHash,
	// 		CreateTime:               timestamppb.New(e.CreateTime),
	// 		UpdateTime:               timestamppb.New(e.UpdateTime),
	// 	})
	// }
	return &pb.GetItemResponse{Item: &pb.Item{
		Id:                       item.ID.String(),
		Name:                     item.Name,
		Description:              item.Description,
		Price:                    item.Price,
		Availability:             int32(item.Availability),
		BusinessFk:               item.BusinessFk.String(),
		BusinessItemCategoryFk:   item.BusinessItemCategoryFk.String(),
		HighQualityPhoto:         item.HighQualityPhoto,
		HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
		LowQualityPhoto:          item.LowQualityPhoto,
		LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
		Thumbnail:                item.Thumbnail,
		ThumbnailBlurHash:        item.ThumbnailBlurHash,
		Cursor:                   item.Cursor,
		CreateTime:               timestamppb.New(item.CreateTime),
		UpdateTime:               timestamppb.New(item.UpdateTime),
	}}, nil
}

func (m *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	var businessItemCategoryFk uuid.UUID
	if req.BusinessItemCategoryFk != "" {
		businessItemCategoryFk = uuid.MustParse(req.BusinessItemCategoryFk)
	}
	item, err := m.itemService.UpdateItem(&dto.UpdateItemRequest{ItemFk: uuid.MustParse(req.ItemFk), Name: req.Name, Description: req.Description, Price: float32(req.Price), HighQualityPhotoObject: req.HighQualityPhotoObject, HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash, LowQualityPhotoObject: req.LowQualityPhotoObject, LowQualityPhotoBlurHash: req.LowQualityPhotoBlurHash, ThumbnailObject: req.ThumbnailObject, ThumbnailBlurHash: req.ThumbnailBlurHash, Availability: req.Availability, Status: req.Status.String(), BusinessItemCategoryFk: businessItemCategoryFk, Metadata: &md})
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
		case "record not found":
			st = status.New(codes.NotFound, "Item not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.UpdateItemResponse{Item: &pb.Item{
		Id:                       item.ID.String(),
		Name:                     item.Name,
		Description:              item.Description,
		Price:                    item.Price,
		Availability:             int32(item.Availability),
		BusinessFk:               item.BusinessFk.String(),
		BusinessItemCategoryFk:   item.BusinessItemCategoryFk.String(),
		HighQualityPhoto:         item.HighQualityPhoto,
		HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
		LowQualityPhoto:          item.LowQualityPhoto,
		LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
		Thumbnail:                item.Thumbnail,
		ThumbnailBlurHash:        item.ThumbnailBlurHash,
		Cursor:                   item.Cursor,
		Status:                   *ut.ParseItemStatusType(&item.Status),
		CreateTime:               timestamppb.New(item.CreateTime),
		UpdateTime:               timestamppb.New(item.UpdateTime),
	}}, nil
}

func (m *ItemServer) SearchItem(ctx context.Context, req *pb.SearchItemRequest) (*pb.SearchItemResponse, error) {
	var st *status.Status
	response, err := m.itemService.SearchItem(req.Name, req.ProvinceFk, req.MunicipalityFk, int64(req.NextPage), req.SearchMunicipalityType.String())
	if err != nil {
		switch err.Error() {
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	itemsResponse := make([]*pb.SearchItem, 0, len(*response.Items))
	for _, e := range *response.Items {
		itemsResponse = append(itemsResponse, &pb.SearchItem{
			Id:                e.ID.String(),
			Name:              e.Name,
			Thumbnail:         e.Thumbnail,
			ThumbnailBlurHash: e.ThumbnailBlurHash,
			Price:             e.Price,
			Cursor:            int32(e.Cursor),
			Status:            *ut.ParseItemStatusType(&e.Status),
		})
	}
	return &pb.SearchItemResponse{Items: itemsResponse, NextPage: response.NextPage, SearchMunicipalityType: *ut.ParseSearchMunicipalityType(response.SearchMunicipalityType)}, nil
}

func (m *ItemServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	err := m.itemService.DeleteItem(&dto.DeleteItemRequest{ItemFk: uuid.MustParse(req.ItemFk), Metadata: &md})
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
	return &gp.Empty{}, nil
}

func (m *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	var st *status.Status
	res, err := m.itemService.CreateItem(&dto.CreateItemRequest{Name: req.Name, Description: req.Description, Price: float32(req.Price), BusinessItemCategoryFk: req.BusinessItemCategoryFk, HighQualityPhotoObject: req.HighQualityPhoto, HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash, LowQualityPhotoObject: req.LowQualityPhoto, LowQualityPhotoBlurHash: req.LowQualityPhotoBlurHash, ThumbnailObject: req.Thumbnail, ThumbnailBlurHash: req.ThumbnailBlurHash, BusinessFk: uuid.MustParse(req.BusinessFk), Metadata: &md})
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
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.CreateItemResponse{Item: &pb.Item{Id: res.ID.String(), Name: res.Name, Description: res.Description, Price: res.Price, Status: *ut.ParseItemStatusType(&res.Status), Availability: int32(res.Availability), BusinessFk: res.BusinessFk.String(), BusinessItemCategoryFk: res.BusinessItemCategoryFk.String(), HighQualityPhoto: res.HighQualityPhoto, HighQualityPhotoBlurHash: res.HighQualityPhotoBlurHash, LowQualityPhoto: res.LowQualityPhoto, LowQualityPhotoBlurHash: res.LowQualityPhotoBlurHash, Thumbnail: res.Thumbnail, ThumbnailBlurHash: res.ThumbnailBlurHash, CreateTime: timestamppb.New(res.CreateTime), UpdateTime: timestamppb.New(res.UpdateTime)}}, nil
}
