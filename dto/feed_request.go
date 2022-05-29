package dto

import (
	"github.com/daniarmas/api_go/utils"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type FeedRequest struct {
	ProvinceId             string
	MunicipalityId         string
	HomeDelivery           bool
	ToPickUp               bool
	NextPage               int32
	Location               ewkb.Point
	SearchMunicipalityType string
	Metadata               *utils.ClientMetadata
}
