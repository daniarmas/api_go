package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type CreateOrderRequest struct {
	CartItems   *[]uuid.UUID
	Status         string
	DeliveryType   string
	ResidenceType  string
	BuildingNumber string
	HouseNumber    string
	Coordinates    ewkb.Point
	DeliveryDate   time.Time
	Metadata       *metadata.MD
}
