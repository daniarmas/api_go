package dto

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type GetAddressInfoRequest struct {
	Coordinates ewkb.Point
	Metadata    *metadata.MD
}
