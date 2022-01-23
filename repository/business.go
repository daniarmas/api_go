package repository

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type BusinessQuery interface {
	Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool) (*[]models.Business, error)
	GetBusiness(tx *gorm.DB, coordinates ewkb.Point, id string) (*models.Business, error)
}

type businessQuery struct{}

func (b *businessQuery) Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool) (*[]models.Business, error) {
	var businessResult *[]models.Business
	point := fmt.Sprintf("'POINT(%v %v)'", coordinates.Point.Coords()[1], coordinates.Point.Coords()[0])
	var where string
	if municipalityNotEqual {
		where = fmt.Sprintf("WHERE cursor > %v AND province_fk = '%v' AND municipality_fk != '%v'", cursor, provinceFk, municipalityFk)
	} else {
		where = fmt.Sprintf("WHERE cursor > %v AND province_fk = '%v' AND municipality_fk = '%v'", cursor, provinceFk, municipalityFk)
	}
	query := fmt.Sprintf("SELECT id, name, address, high_quality_photo, high_quality_photo_blurhash, low_quality_photo, low_quality_photo_blurhash, delivery_price, is_open, home_delivery, to_pick_up, cursor, ST_AsEWKB(business.coordinates) AS coordinates, ST_AsEWKB(business.polygon) AS polygon, ST_Contains(business.polygon, ST_GeomFromText(%v, 4326)) as is_in_range FROM business %v ORDER BY cursor asc LIMIT 6;", point, where)
	err := tx.Raw(query).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	}
	return businessResult, nil
}

func (b *businessQuery) GetBusiness(tx *gorm.DB, coordinates ewkb.Point, id string) (*models.Business, error) {
	var businessResult *models.Business
	point := fmt.Sprintf("'POINT(%v %v)'", coordinates.Point.Coords()[1], coordinates.Point.Coords()[0])
	query := fmt.Sprintf("SELECT business.id, business.name, business.address, business.high_quality_photo, business.high_quality_photo_blurhash, business.low_quality_photo, business.low_quality_photo_blurhash, business.delivery_price, business.is_open, business.home_delivery, business.to_pick_up, business.cursor, ST_AsEWKB(business.coordinates) AS coordinates, ST_AsEWKB(business.polygon) AS polygon, ST_Contains(business.polygon, ST_GeomFromText(%v, 4326)) as is_in_range FROM business WHERE business.id = '%v' ORDER BY business.cursor asc LIMIT 1;", point, id)
	err := tx.Raw(query).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	} else if businessResult == nil {
		return nil, errors.New("business not found")
	}
	return businessResult, nil
}
