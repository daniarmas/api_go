package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UpdateItemRequest struct {
	ItemId                   uuid.UUID
	Name                     string
	Description              string
	Price                    string
	HighQualityPhoto         string
	HighQualityPhotoBlurHash string
	LowQualityPhoto          string
	LowQualityPhotoBlurHash  string
	Thumbnail                string
	ThumbnailBlurHash        string
	Availability             int64
	Status                   string
	BusinessColletionId      uuid.UUID
	BusinessId               uuid.UUID
	Metadata                 *metadata.MD
}
