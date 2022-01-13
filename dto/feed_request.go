package dto

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type FeedRequest struct {
	ProvinceFk             string
	MunicipalityFk         string
	NextPage               int32
	Location               ewkb.Point
	SearchMunicipalityType string
}
