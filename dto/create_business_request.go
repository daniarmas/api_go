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
	HighQualityPhoto         string
	HighQualityPhotoBlurHash string
	LowQualityPhoto          string
	LowQualityPhotoBlurHash  string
	Thumbnail                string
	ThumbnailBlurHash        string
	Municipalities           []string
	DeliveryPrice            string
	Coordinates              ewkb.Point
	TimeMarginOrderMonth     int32
	TimeMarginOrderDay       int32
	TimeMarginOrderHour      int32
	TimeMarginOrderMinute    int32
	ToPickUp                 bool
	HomeDelivery             bool
	BusinessBrandId          string
	ProvinceId               string
	MunicipalityId           string
}
