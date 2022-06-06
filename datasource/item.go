package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemDatasource interface {
	GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error)
	GetItemWithLocation(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error)
	ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error)
	ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error)
	CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceid string, municipalityid string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error)
	UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error)
	UpdateItems(tx *gorm.DB, data *[]models.Item) (*[]models.Item, error)
	DeleteItem(tx *gorm.DB, where *models.Item) error
}

type itemDatasource struct{}

func (v *itemDatasource) DeleteItem(tx *gorm.DB, where *models.Item) error {
	var itemResult *[]models.Item
	result := tx.Where(where).Delete(&itemResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *itemDatasource) CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *itemDatasource) ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error) {
	var res []models.Item
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Limit(11).Select(*selectFields).Where(where).Where("create_time < ?", cursor).Order("create_time desc").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *itemDatasource) ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error) {
	var res []models.Item
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where("id IN ? ", ids).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *itemDatasource) GetItemWithLocation(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error) {
	var res *models.ItemBusiness
	// p := fmt.Sprintf("'POINT(%v %v)'", point.Point.Coords()[1], point.Point.Coords()[0])
	result := tx.Model(&models.Item{}).Select("item.id, item.name, item.business_collection_id, item.business_id, item.description, item.price, item.availability, item.business_id, item.high_quality_photo, item.high_quality_photo_blurhash, item.low_quality_photo, item.low_quality_photo_blurhash, item.thumbnail, item.thumbnail_blurhash, item.create_time, item.update_time, item.cursor").Where("item.id = ?", id).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (i *itemDatasource) GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error) {
	var res *models.Item
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (i *itemDatasource) SearchItem(tx *gorm.DB, name string, provinceid string, municipalityid string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error) {
	var items []models.Item
	var result *gorm.DB
	where := "%" + name + "%"
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	if municipalityNotEqual {
		result = tx.Limit(int(limit+1)).Select(*selectFields).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_id != ? AND province_id = ? AND cursor > ?", where, municipalityid, provinceid, cursor).Order("cursor asc").Find(&items)
	} else {
		result = tx.Limit(int(limit+1)).Select(*selectFields).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_id = ? AND province_id = ? AND cursor > ?", where, municipalityid, provinceid, cursor).Order("cursor asc").Find(&items)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (v *itemDatasource) UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error) {
	result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}

func (v *itemDatasource) UpdateItems(tx *gorm.DB, data *[]models.Item) (*[]models.Item, error) {
	result := tx.Clauses(clause.Returning{}).Updates(data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}
