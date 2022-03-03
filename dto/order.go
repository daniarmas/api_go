package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type Order struct {
	ID                   uuid.UUID
	BusinessName         string
	ItemQuantity         int
	Status               string
	DeliveryType         string
	ResidenceType        string
	Price                float64
	BuildingNumber       string
	HouseNumber          string
	BusinessFk           uuid.UUID
	Coordinates          ewkb.Point
	UserFk               uuid.UUID
	AuthorizationTokenFk uuid.UUID
	DeliveryDate         time.Time
	CreateTime           time.Time
	UpdateTime           time.Time
}
