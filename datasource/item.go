package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemDatasource interface {
	GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error)
	ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error)
	ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error)
	CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error)
	SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]models.Item, error)
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
	if fields != nil {
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
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where("id IN ? ", ids).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *itemDatasource) GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error) {
	var res *models.Item
	selectFields := &[]string{"*"}
	if fields != nil {
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

func (i *itemDatasource) SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error) {
	var items []models.Item
	var result *gorm.DB
	where := "%" + name + "%"
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	if municipalityNotEqual {
		result = tx.Limit(int(limit+1)).Select(*selectFields).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_id != ? AND province_id = ? AND cursor > ?", where, municipalityId, provinceId, cursor).Order("cursor asc").Find(&items)
	} else {
		result = tx.Limit(int(limit+1)).Select(*selectFields).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_id = ? AND province_id = ? AND cursor > ?", where, municipalityId, provinceId, cursor).Order("cursor asc").Find(&items)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (i *itemDatasource) SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]models.Item, error) {
	var items []models.Item
	var result *gorm.DB
	where := "%" + name + "%"
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result = tx.Limit(10).Select(*selectFields).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND cursor > ? AND business_id = ?", where, cursor, businessId).Order("cursor asc").Find(&items)
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
