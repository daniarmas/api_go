package dto

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type FeedRequest struct {
	ProvinceFk             string
	MunicipalityFk         string
	HomeDelivery           bool
	ToPickUp               bool
	NextPage               int32
	Location               ewkb.Point
	SearchMunicipalityType string
	Metadata               *metadata.MD
}
