package dto

import "github.com/twpayne/go-geom/encoding/ewkb"

type GetBusinessRequest struct {
	Id          string
	Coordinates ewkb.Point
}
