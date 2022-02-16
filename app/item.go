package app

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	itemPhotos := make([]*pb.ItemPhoto, 0, len(item.ItemPhoto))
	for _, e := range item.ItemPhoto {
		itemPhotos = append(itemPhotos, &pb.ItemPhoto{
			Id:                       e.ID.String(),
			ItemFk:                   e.ItemFk.String(),
			HighQualityPhoto:         e.HighQualityPhoto,
			HighQualityPhotoBlurHash: e.HighQualityPhotoBlurHash,
			LowQualityPhoto:          e.LowQualityPhoto,
			LowQualityPhotoBlurHash:  e.LowQualityPhotoBlurHash,
			Thumbnail:                e.Thumbnail,
			ThumbnailBlurHash:        e.ThumbnailBlurHash,
			CreateTime:               timestamppb.New(e.CreateTime),
			UpdateTime:               timestamppb.New(e.UpdateTime),
		})
	}
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
		Photos:                   itemPhotos,
		IsInRange:                item.IsInRange,
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
			Status:            *utils.ParseItemStatusType(&e.Status),
		})
	}
	return &pb.SearchItemResponse{Items: itemsResponse, NextPage: response.NextPage, SearchMunicipalityType: *utils.ParseSearchMunicipalityType(response.SearchMunicipalityType)}, nil
}
