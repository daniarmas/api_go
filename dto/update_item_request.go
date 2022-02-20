package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UpdateItemRequest struct {
	ItemFk                   uuid.UUID
	Name                     string
	Description              string
	Price                    float32
	HighQualityPhotoObject   string
	HighQualityPhotoBlurHash string
	LowQualityPhotoObject    string
	LowQualityPhotoBlurHash  string
	ThumbnailObject          string
	ThumbnailBlurHash        string
	Availability             int64
	Status                   string
	BusinessItemCategoryFk   uuid.UUID
	BusinessFk               uuid.UUID
	Metadata                 *metadata.MD
}
