package dto

type GetPresignedPutObjectResponse struct {
	LowQualityPhotoPresignedPutUrl       string
	HighQualityPhotoPresignedPutUrl      string
	ThumbnailQualityPhotoPresignedPutUrl string
}
