package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
)

func (m *ItemServer) ListItem(ctx context.Context, req *pb.ListItemRequest) (*pb.ListItemResponse, error) {
	items, err := m.itemService.ListItem()
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*pb.Item, 0, len(items))
	for _, item := range items {
		itemsResponse = append(itemsResponse, &pb.Item{
			Id:                       item.ID,
			Name:                     item.Name,
			Description:              item.Description,
			Price:                    item.Price,
			Availability:             int32(item.Availability),
			BusinessFk:               item.BusinessFk,
			BusinessItemCategoryFk:   item.BusinessItemCategoryFk,
			HighQualityPhoto:         item.HighQualityPhoto,
			HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
			LowQualityPhoto:          item.LowQualityPhoto,
			LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
			Thumbnail:                item.Thumbnail,
			ThumbnailBlurHash:        item.ThumbnailBlurHash,
		})
	}
	return &pb.ListItemResponse{Items: itemsResponse}, nil
}

func (m *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := m.itemService.GetItem(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetItemResponse{Item: &pb.Item{
		Id:                       item.ID,
		Name:                     item.Name,
		Description:              item.Description,
		Price:                    item.Price,
		Availability:             int32(item.Availability),
		BusinessFk:               item.BusinessFk,
		BusinessItemCategoryFk:   item.BusinessItemCategoryFk,
		HighQualityPhoto:         item.HighQualityPhoto,
		HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
		LowQualityPhoto:          item.LowQualityPhoto,
		LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
		Thumbnail:                item.Thumbnail,
		ThumbnailBlurHash:        item.ThumbnailBlurHash,
	}}, nil
}
