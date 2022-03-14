package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type BusinessDatasource interface {
	Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error)
	GetBusiness(tx *gorm.DB, where *models.Business) (*models.Business, error)
	CreateBusiness(tx *gorm.DB, data *models.Business) (*models.Business, error)
	GetBusinessWithLocation(tx *gorm.DB, where *models.Business) (*models.Business, error)
	GetBusinessProvinceAndMunicipality(tx *gorm.DB, businessFk uuid.UUID) (*dto.GetBusinessProvinceAndMunicipality, error)
}

type businessDatasource struct{}

func (b *businessDatasource) GetBusiness(tx *gorm.DB, where *models.Business) (*models.Business, error) {
	var businessResult *models.Business
	result := tx.Where(where).Take(&businessResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessResult, nil
}

func (b *businessDatasource) CreateBusiness(tx *gorm.DB, data *models.Business) (*models.Business, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var response models.Business
	result := tx.Raw(`INSERT INTO "business" ("id", "name", "description", "address", "phone", "email", "high_quality_photo", "high_quality_photo_object", "high_quality_photo_blurhash", "low_quality_photo", "low_quality_photo_object", "low_quality_photo_blurhash", "thumbnail", "thumbnail_object", "thumbnail_blurhash", "is_open", "time_margin_order_month", "time_margin_order_day", "time_margin_order_hour", "time_margin_order_minute", "delivery_price", "to_pick_up", "home_delivery", "coordinates", "province_fk", "municipality_fk", "business_brand_fk",  "create_time", "update_time") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "id", "name", "description", "address", "phone", "email", "high_quality_photo", "high_quality_photo_object", "high_quality_photo_blurhash", "low_quality_photo", "low_quality_photo_object", "low_quality_photo_blurhash", "thumbnail", "thumbnail_object", "thumbnail_blurhash", "is_open", "time_margin_order_month", "time_margin_order_day", "time_margin_order_hour", "time_margin_order_minute", "delivery_price", "to_pick_up", "home_delivery", "coordinates", "province_fk", "municipality_fk", "business_brand_fk",  "create_time", "update_time"`, uuid.New().String(), data.Name, data.Description, data.Address, data.Phone, point, data.Email, data.HighQualityPhoto, data.HighQualityPhotoObject, data.HighQualityPhotoBlurHash, data.LowQualityPhoto, data.LowQualityPhotoObject, data.LowQualityPhotoBlurHash, data.Thumbnail, data.ThumbnailObject, data.ThumbnailBlurHash, data.IsOpen, data.TimeMarginOrderMonth, data.TimeMarginOrderDay, data.TimeMarginOrderHour, data.TimeMarginOrderHour, data.DeliveryPrice, data.ToPickUp, data.HomeDelivery, point, data.ProvinceFk, data.MunicipalityFk, data.BusinessBrandFk, time, time).Scan(&response)
	if result.Error != nil {
		return nil, result.Error
	}
	return &response, nil
}

func (b *businessDatasource) GetBusinessProvinceAndMunicipality(tx *gorm.DB, businessFk uuid.UUID) (*dto.GetBusinessProvinceAndMunicipality, error) {
	var res dto.GetBusinessProvinceAndMunicipality
	result := tx.Model(&models.Business{}).Limit(1).Select("province_fk", "municipality_fk").Where("id = ?", businessFk.String()).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (b *businessDatasource) Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error) {
	var businessResult *[]models.Business
	var delivery string
	if homeDelivery {
		delivery = "business.home_delivery = true"
	} else {
		delivery = "business.to_pick_up = true"
	}
	var where string
	if municipalityNotEqual {
		where = fmt.Sprintf("WHERE cursor > %v AND province_fk = '%v' AND municipality_fk != '%v' AND %v", cursor, provinceFk, municipalityFk, delivery)
	} else {
		where = fmt.Sprintf("WHERE cursor > %v AND province_fk = '%v' AND municipality_fk = '%v' AND %v", cursor, provinceFk, municipalityFk, delivery)
	}
	query := fmt.Sprintf("SELECT id, name, address, high_quality_photo, high_quality_photo_blurhash, low_quality_photo, low_quality_photo_blurhash, delivery_price, is_open, home_delivery, to_pick_up, cursor FROM business %v AND status = 'BusinessAvailable' ORDER BY cursor asc LIMIT 6;", where)
	err := tx.Raw(query).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	}
	return businessResult, nil
}

func (b *businessDatasource) GetBusinessWithLocation(tx *gorm.DB, where *models.Business) (*models.Business, error) {
	var businessResult *models.Business
	// point := fmt.Sprintf("'POINT(%v %v)'", where.Coordinates.Point.Coords()[1], where.Coordinates.Point.Coords()[0])
	distance := fmt.Sprintf(`ST_Distance("coordinates", ST_GeomFromText('POINT(%v %v)', 4326)) AS "distance"`, where.Coordinates.Point.Coords()[1], where.Coordinates.Point.Coords()[0])
	query := fmt.Sprintf("SELECT business.id, business.name, business.phone, business.description, business.email, business.address, business.high_quality_photo, business.high_quality_photo_blurhash, business.time_margin_order_month, business.time_margin_order_day, business.time_margin_order_hour, business.time_margin_order_minute, business.low_quality_photo, business.low_quality_photo_blurhash, business.delivery_price, business.is_open, business.home_delivery, business.to_pick_up, business.cursor, ST_AsEWKB(business.coordinates) AS coordinates, %v FROM business WHERE business.id = '%v' ORDER BY business.cursor asc LIMIT 1;", distance, where.ID)
	err := tx.Raw(query).Scan(&businessResult).Error
	if err != nil {
		return nil, err
	} else if businessResult == nil {
		return nil, errors.New("business not found")
	}
	return businessResult, nil
}
