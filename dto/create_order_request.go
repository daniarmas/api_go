package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/metadata"
)

type CreateOrderRequest struct {
	CartItems     *[]uuid.UUID
	OrderType     string
	ResidenceType string
	Number        string
	Address       string
	Instructions  string
	Coordinates   ewkb.Point
	OrderTime     time.Time
	Metadata      *metadata.MD
}
