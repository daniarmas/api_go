package dto

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type CreateBusinessRequest struct {
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
	Municipalities           []string
	IsOpen                   bool
	DeliveryPrice            float64
	Coordinates              ewkb.Point
	TimeMarginOrderMonth     int32
	TimeMarginOrderDay       int32
	TimeMarginOrderHour      int32
	TimeMarginOrderMinute    int32
	ToPickUp                 bool
	HomeDelivery             bool
	BusinessBrandFk          string
	ProvinceFk               string
	MunicipalityFk           string
}
