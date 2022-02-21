package dto

import (
	"time"

	"github.com/twpayne/go-geom/encoding/ewkb"
)

type CreateOrderRequest struct {
	OrderedItems   *[]string
	Status         string
	DeliveryType   string
	ResidenceType  string
	BuildingNumber string
	HouseNumber    string
	Coordinates    ewkb.Point
	DeliveryDate   time.Time
}
