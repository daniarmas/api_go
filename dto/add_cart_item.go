package dto

import (
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type AddCartItem struct {
	ItemId         string
	Quantity       int32
	Location       ewkb.Point
	MunicipalityId uuid.UUID
	Metadata       *metadata.MD
}
