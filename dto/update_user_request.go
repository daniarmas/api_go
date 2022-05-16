package dto

import "google.golang.org/grpc/metadata"

type UpdateUserRequest struct {
	Id                       string
	Email                    string
	FullName                 string
	Thumbnail                string
	ThumbnailBlurHash        string
	HighQualityPhoto         string
	HighQualityPhotoBlurHash string
	LowQualityPhoto          string
	LowQualityPhotoBlurHash  string
	Code                     string
	Metadata                 *metadata.MD
}
