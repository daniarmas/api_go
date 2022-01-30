package dto

import "google.golang.org/grpc/metadata"

type UpdateUserRequest struct {
	Id                       string
	Email                    string
	Alias                    string
	FullName                 string
	ThumbnailObject          string
	ThumbnailBlurHash        string
	HighQualityPhotoObject   string
	HighQualityPhotoBlurHash string
	LowQualityPhotoObject    string
	LowQualityPhotoBlurHash  string
	Code                     string
	Metadata                 *metadata.MD
}
