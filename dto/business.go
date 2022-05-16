package dto

import (
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type Business struct {
	ID                       uuid.UUID
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
	IsOpen                   bool
	IsInRange                bool
	DeliveryPrice            string
	Coordinates              ewkb.Point
	Polygon                  ewkb.Polygon
	TimeMarginOrderMonth     int32
	TimeMarginOrderDay       int32
	TimeMarginOrderHour      int32
	TimeMarginOrderMinute    int32
	ToPickUp                 bool
	HomeDelivery             bool
	BusinessBrandId          uuid.UUID
	ProvinceId               uuid.UUID
	MunicipalityId           uuid.UUID
	Distance                 float32
	Status                   string
	Cursor                   int32
}
