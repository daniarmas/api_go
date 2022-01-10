package repository

import (
	"fmt"

	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type BusinessQuery interface {
	ListBusiness(tx *gorm.DB, where *datastruct.Business) (*[]datastruct.Business, error)
}

type businessQuery struct{}

func (b *businessQuery) ListBusiness(tx *gorm.DB, where *datastruct.Business) (*[]datastruct.Business, error) {
	var businessResult *[]datastruct.Business
	point := fmt.Sprintf("POINT(%v %v)", where.Coordinates.Point.Coords()[1], where.Coordinates.Point.Coords()[0])
	err := tx.Raw(`SELECT id, name, description, address, phone, email, high_quality_photo, high_quality_photo_blurhash, low_quality_photo, low_quality_photo_blurhash, thumbnail, thumbnail_blurhash, delivery_price, is_open, lead_day_time, lead_hours_time, lead_minutes_time, home_delivery, to_pick_up, business_brand_fk, province_fk, municipality_fk, ST_AsEWKB(business.coordinates) AS coordinates, ST_AsEWKB(business.polygon) AS polygon, ST_Contains(business.polygon, ST_GeomFromText(?, 4326)) as is_in_range FROM business;`, point).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	}
	return businessResult, nil
}
