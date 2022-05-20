package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type CreateItemRequest struct {
	Name                     string
	Description              string
	Price                    string
	BusinessCollectionId     string
	HighQualityPhoto         string
	HighQualityPhotoBlurHash string
	LowQualityPhoto          string
	LowQualityPhotoBlurHash  string
	Thumbnail                string
	ThumbnailBlurHash        string
	BusinessId               *uuid.UUID
	Metadata                 *metadata.MD
}
