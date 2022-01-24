package dto

import "github.com/twpayne/go-geom/encoding/ewkb"

type GetItemRequest struct {
	Id       string
	Location ewkb.Point
}
