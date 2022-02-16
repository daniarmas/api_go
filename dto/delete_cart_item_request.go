package dto

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type DeleteCartItemRequest struct {
	CartItemFk string
	Location   ewkb.Point
	Metadata   *metadata.MD
}
