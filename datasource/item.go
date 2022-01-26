package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type ItemDatasource interface {
	GetItem(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error)
	ListItem(tx *gorm.DB, where *models.Item) (*[]models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error)
}

type itemDatasource struct{}

func (i *itemDatasource) ListItem(tx *gorm.DB, where *models.Item) (*[]models.Item, error) {
	var items []models.Item
	result := tx.Limit(11).Where("business_fk = ? AND business_item_category_fk = ? AND cursor > ?", where.BusinessFk, where.BusinessItemCategoryFk, where.Cursor).Order("cursor asc").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (i *itemDatasource) GetItem(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error) {
	var item *models.ItemBusiness
	p := fmt.Sprintf("'POINT(%v %v)'", point.Point.Coords()[1], point.Point.Coords()[0])
	s := fmt.Sprintf("item.id, item.name, item.description, item.price, item.status, item.availability, item.business_fk, item.business_item_category_fk, item.high_quality_photo, item.high_quality_photo_blurhash, item.low_quality_photo, item.low_quality_photo_blurhash, item.thumbnail, item.thumbnail_blurhash, item.create_time, item.update_time, item.cursor, ST_Contains(business.polygon, ST_GeomFromText(%s, 4326)) as is_in_range", p)
	result := tx.Model(&models.Item{}).Preload("ItemPhoto").Select(s).Joins("left join business on business.id = item.business_fk").Where("item.id = ?", id).Take(&item)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return item, nil
}

func (i *itemDatasource) SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error) {
	var items []models.Item
	var result *gorm.DB
	where := "%" + name + "%"
	if municipalityNotEqual {
		result = tx.Limit(int(limit+1)).Select("id, name, price, thumbnail, thumbnail_blurhash, cursor, status").Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_fk != ? AND province_fk = ? AND cursor > ?", where, municipalityFk, provinceFk, cursor).Order("cursor asc").Find(&items)
	} else {
		result = tx.Limit(int(limit+1)).Select("id, name, price, thumbnail, thumbnail_blurhash, cursor, status").Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_fk = ? AND province_fk = ? AND cursor > ?", where, municipalityFk, provinceFk, cursor).Order("cursor asc").Find(&items)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}
