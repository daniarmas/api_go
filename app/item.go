package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *ItemServer) ListItem(ctx context.Context, req *pb.ListItemRequest) (*pb.ListItemResponse, error) {
	listItemsResponse, err := m.itemService.ListItem(&dto.ListItemRequest{BusinessFk: req.BusinessFk, BusinessItemCategoryFk: req.ItemCategoryFk, NextPage: req.NextPage})
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
			CreateTime:               item.CreateTime.String(),
			UpdateTime:               item.UpdateTime.String(),
			Cursor:                   int32(item.Cursor),
		})
	}
	return &pb.ListItemResponse{Items: itemsResponse, NextPage: listItemsResponse.NextPage}, nil
}

func (m *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := m.itemService.GetItem(req.Id)
	if err != nil {
		return nil, err
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
