package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessDatasource interface {
	Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceId string, municipalityId string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]entity.Business, error)
	GetBusiness(tx *gorm.DB, where *entity.Business, fields *[]string) (*entity.Business, error)
	CreateBusiness(tx *gorm.DB, data *entity.Business) (*entity.Business, error)
	UpdateBusiness(tx *gorm.DB, data *entity.Business, where *entity.Business) (*entity.Business, error)
	UpdateBusinessCoordinate(tx *gorm.DB, data *entity.Business, where *entity.Business) error
	GetBusinessWithDistance(tx *gorm.DB, where *entity.Business) (*entity.Business, error)
	BusinessIsInRange(tx *gorm.DB, coordinates ewkb.Point, businessId *uuid.UUID) (*bool, error)
}

type businessDatasource struct{}

func (b *businessDatasource) BusinessIsInRange(tx *gorm.DB, coordinates ewkb.Point, businessId *uuid.UUID) (*bool, error) {
	type IsInRange struct {
		IsInRange bool `gorm:"column:is_in_range;not null"`
	}
	var res *IsInRange
	p := fmt.Sprintf("ST_Contains(business.polygon, ST_GeomFromText('POINT(%v %v)', 4326)) as is_in_range", coordinates.Point.Coords()[1], coordinates.Point.Coords()[0])
	result := tx.Model(&entity.Business{}).Select(p).Where("id = ?", businessId).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res.IsInRange, nil
}

func (b *businessDatasource) GetBusiness(tx *gorm.DB, where *entity.Business, fields *[]string) (*entity.Business, error) {
	var res *entity.Business
	var selectField *[]string
	if fields == nil {
		selectField = &[]string{"id", "name", "address", "high_quality_photo", "low_quality_photo", "thumbnail", "blurhash", "time_margin_order_month", "time_margin_order_day", "time_margin_order_hour", "time_margin_order_minute", "delivery_price_cup", "to_pick_up", "home_delivery", "ST_AsEWKB(coordinates) AS coordinates", "province_id", "municipality_id", "business_brand_id", "enabled_flag", "create_time", "update_time", "cursor"}
	} else {
		selectField = fields
	}
	result := tx.Select(*selectField).Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (b *businessDatasource) CreateBusiness(tx *gorm.DB, data *entity.Business) (*entity.Business, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var res entity.Business
	var countRes []entity.Business
	number := tx.Select("id").Where("municipality_id = ?", data.MunicipalityId).Find(&countRes)
	result := tx.Raw(`INSERT INTO "business" ("id", "name", "address", "high_quality_photo", "low_quality_photo", "thumbnail", "blurhash", "time_margin_order_month", "time_margin_order_day", "time_margin_order_hour", "time_margin_order_minute", "delivery_price_cup", "to_pick_up", "home_delivery", "coordinates", "province_id", "municipality_id", "business_brand_id",  "create_time", "update_time", "cursor") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?) RETURNING "id", "name", "description", "address", "high_quality_photo", "low_quality_photo", "thumbnail", "blurhash", "time_margin_order_month", "time_margin_order_day", "time_margin_order_hour", "time_margin_order_minute", "delivery_price_cup", "to_pick_up", "home_delivery", ST_AsEWKB(coordinates) AS coordinates, "province_id", "municipality_id", "business_brand_id",  "create_time", "update_time", "cursor"`, uuid.New().String(), data.Name, data.Address, data.HighQualityPhoto, data.LowQualityPhoto, data.Thumbnail, data.BlurHash, data.TimeMarginOrderMonth, data.TimeMarginOrderDay, data.TimeMarginOrderHour, data.TimeMarginOrderMinute, data.DeliveryPriceCup, data.ToPickUp, data.HomeDelivery, point, data.ProvinceId, data.MunicipalityId, data.BusinessBrandId, time, time, number.RowsAffected+1).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (b *businessDatasource) UpdateBusiness(tx *gorm.DB, data *entity.Business, where *entity.Business) (*entity.Business, error) {
	result := tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "name"}, {Name: "description"}, {Name: "address"}, {Name: "high_quality_photo"}, {Name: "high_quality_photo_blurhash"}, {Name: "low_quality_photo"}, {Name: "low_quality_photo_blurhash"}, {Name: "thumbnail"}, {Name: "thumbnail_blurhash"}, {Name: "time_margin_order_month"}, {Name: "time_margin_order_day"}, {Name: "time_margin_order_hour"}, {Name: "time_margin_order_minute"}, {Name: "delivery_price_cup"}, {Name: "to_pick_up"}, {Name: "home_delivery"}, {Name: "home_delivery"}, {Name: "province_id"}, {Name: "municipality_id"}, {Name: "business_brand_id"}, {Name: "create_time"}, {Name: "update_time"}}}).Where(where).Updates(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}

func (b *businessDatasource) UpdateBusinessCoordinate(tx *gorm.DB, data *entity.Business, where *entity.Business) error {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var response entity.Business
	result := tx.Raw(`UPDATE "business" SET coordinates = ST_GeomFromText(?, 4326), update_time = ? WHERE id = ?`, point, time, where.ID).Scan(&response)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (b *businessDatasource) Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceId string, municipalityId string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]entity.Business, error) {
	var businessResult *[]entity.Business
	var delivery string
	if homeDelivery {
		delivery = "business.home_delivery = true"
	} else {
		delivery = "business.to_pick_up = true"
	}
	var where string
	if municipalityNotEqual {
		where = fmt.Sprintf("WHERE cursor > %v AND province_id = '%v' AND municipality_id != '%v' AND %v", cursor, provinceId, municipalityId, delivery)
	} else {
		where = fmt.Sprintf("WHERE cursor > %v AND province_id = '%v' AND municipality_id = '%v' AND %v", cursor, provinceId, municipalityId, delivery)
	}
	query := fmt.Sprintf("SELECT id, name, address, high_quality_photo, low_quality_photo, blurhash, delivery_price_cup, home_delivery, to_pick_up, business_brand_id, province_id, municipality_id, cursor FROM business %v ORDER BY cursor asc LIMIT 6;", where)
	err := tx.Raw(query).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	}
	return businessResult, nil
}

func (b *businessDatasource) GetBusinessWithDistance(tx *gorm.DB, where *entity.Business) (*entity.Business, error) {
	var businessResult *entity.Business
	distance := fmt.Sprintf(`ST_Distance("coordinates", ST_GeomFromText('POINT(%v %v)', 4326)) AS "distance"`, where.Coordinates.Point.Coords()[1], where.Coordinates.Point.Coords()[0])
	query := fmt.Sprintf("SELECT business.id, business.name, business_category.name as business_category, business.business_brand_id, business.municipality_id, business.province_id, business.thumbnail, business.blurhash, business.address, business.high_quality_photo, business.time_margin_order_month, business.time_margin_order_day, business.time_margin_order_hour, business.time_margin_order_minute, business.low_quality_photo, business.delivery_price_cup, business.home_delivery, business.to_pick_up, business.cursor, ST_AsEWKB(business.coordinates) AS coordinates, %v FROM business INNER JOIN business_category ON business.business_category_id=business_category.id WHERE business.id = '%v' ORDER BY business.cursor asc LIMIT 1;", distance, where.ID)
	err := tx.Raw(query).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	} else if businessResult == nil {
		return nil, errors.New("record not found")
	}
	return businessResult, nil
}