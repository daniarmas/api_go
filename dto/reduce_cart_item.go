package dto

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type ReduceCartItem struct {
	ItemFk   string
	Location ewkb.Point
	Metadata *metadata.MD
}
