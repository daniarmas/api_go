package dto

import "google.golang.org/grpc/metadata"

type GetPresignedPutObjectRequest struct {
	Metadata                    *metadata.MD
	PhotoType                   string
	LowQualityPhotoObject       string
	HighQualityPhotoObject      string
	ThumbnailQualityPhotoObject string
}
