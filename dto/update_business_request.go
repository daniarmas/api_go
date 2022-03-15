package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UpdateBusinessRequest struct {
	Id                       uuid.UUID
	Metadata                 *metadata.MD
	Name                     string
	Description              string
	Address                  string
	Phone                    string
	Email                    string
	HighQualityPhotoObject   string
	HighQualityPhotoBlurHash string
	LowQualityPhotoObject    string
	LowQualityPhotoBlurHash  string
	ThumbnailObject          string
	ThumbnailBlurHash        string
	DeliveryPrice            float64
	TimeMarginOrderMonth     int32
	TimeMarginOrderDay       int32
	TimeMarginOrderHour      int32
	TimeMarginOrderMinute    int32
	ToPickUp                 bool
	HomeDelivery             bool
	ProvinceFk               string
	MunicipalityFk           string
}
