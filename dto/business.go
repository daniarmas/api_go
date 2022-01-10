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
	DeliveryPrice            float32
	Coordinates              ewkb.Point
	Polygon                  ewkb.Polygon
	LeadDayTime              int32
	LeadHoursTime            int32
	LeadMinutesTime          int32
	ToPickUp                 bool
	HomeDelivery             bool
	BusinessBrandFk          uuid.UUID
	ProvinceFk               uuid.UUID
	MunicipalityFk           uuid.UUID
	Distance                 float32
	Status                   string
	Cursor                   int32
}
