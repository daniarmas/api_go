package dto

import (
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type DeleteCartItemRequest struct {
	CartItemId     string
	Location       ewkb.Point
	Metadata       *metadata.MD
	MunicipalityId *uuid.UUID
}
