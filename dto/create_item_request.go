package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type CreateItemRequest struct {
	Name                     string
	Description              string
	Price                    float32
	BusinessItemCategoryFk   string
	HighQualityPhotoObject   string
	HighQualityPhotoBlurHash string
	LowQualityPhotoObject    string
	LowQualityPhotoBlurHash  string
	ThumbnailObject          string
	ThumbnailBlurHash        string
	BusinessFk               uuid.UUID
	Metadata                 *metadata.MD
}
