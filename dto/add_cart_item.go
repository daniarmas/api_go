package dto

import (
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type AddCartItem struct {
	ItemFk         string
	Quantity       int32
	Location       ewkb.Point
	MunicipalityFk uuid.UUID
	Metadata       *metadata.MD
}
