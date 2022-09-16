package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemDatasource interface {
	GetItem(tx *gorm.DB, where *entity.ItemBusiness) (*entity.ItemBusiness, error)
	ListItem(tx *gorm.DB, where *entity.ItemBusiness, cursor time.Time) (*[]entity.ItemBusiness, error)
	ListItemInIds(tx *gorm.DB, ids []uuid.UUID) (*[]entity.ItemBusiness, error)
	CreateItem(tx *gorm.DB, data *entity.ItemBusiness) (*entity.ItemBusiness, error)
	SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64) (*[]entity.ItemBusiness, error)
	SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string) (*[]entity.ItemBusiness, error)
	UpdateItem(tx *gorm.DB, where *entity.ItemBusiness, data *entity.ItemBusiness) (*entity.ItemBusiness, error)
	UpdateItems(tx *gorm.DB, data *[]entity.ItemBusiness) (*[]entity.ItemBusiness, error)
	DeleteItem(tx *gorm.DB, where *entity.ItemBusiness, ids *[]uuid.UUID) (*[]entity.ItemBusiness, error)
}

type itemDatasource struct{}

func (v *itemDatasource) DeleteItem(tx *gorm.DB, where *entity.ItemBusiness, ids *[]uuid.UUID) (*[]entity.ItemBusiness, error) {
	var res *[]entity.ItemBusiness
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *itemDatasource) CreateItem(tx *gorm.DB, data *entity.ItemBusiness) (*entity.ItemBusiness, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *itemDatasource) ListItem(tx *gorm.DB, where *entity.ItemBusiness, cursor time.Time) (*[]entity.ItemBusiness, error) {
	var res []entity.ItemBusiness
	result := tx.Limit(11).Where(where).Where("create_time < ?", cursor).Order("create_time desc").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *itemDatasource) ListItemInIds(tx *gorm.DB, ids []uuid.UUID) (*[]entity.ItemBusiness, error) {
	var res []entity.ItemBusiness
	result := tx.Where("id IN ? ", ids).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *itemDatasource) GetItem(tx *gorm.DB, where *entity.ItemBusiness) (*entity.ItemBusiness, error) {
	var res *entity.ItemBusiness
	result := tx.Model(entity.Item{}).Select("item.id, item.name, business.name as business_name, item.description,  item.province_id, item.municipality_id, item.business_collection_id, item.business_id, item.availability, item.enabled_flag, item.available_flag, business.open_flag, item.profit_usd, item.cost_usd, item.price_usd, item.profit_cup, item.cost_cup, item.price_cup, item.high_quality_photo, item.low_quality_photo, item.thumbnail, item.blurhash, item.cursor, item.create_time, item.update_time").Joins("left join business on business.id = item.business_id").Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
	// ID                   *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	// Name                 string         `gorm:"column:name;not null"`
	// BusinessName         string         `gorm:"column:business_name;not null"`
	// Description          string         `gorm:"column:description"`
	// PriceCup             string         `gorm:"column:price_cup;not null"`
	// CostCup              string         `gorm:"column:cost_cup"`
	// ProfitCup            string         `gorm:"column:profit_cup"`
	// PriceUsd             string         `gorm:"column:price_usd"`
	// CostUsd              string         `gorm:"column:cost_usd"`
	// ProfitUsd            string         `gorm:"column:profit_usd"`
	// BusinessOpenFlag     bool           `gorm:"column:open_flag;not null"`
	// AvailableFlag        bool           `gorm:"column:available_flag;not null"`
	// EnabledFlag          bool           `gorm:"column:enabled_flag;not null"`
	// Availability         int64          `gorm:"column:availability;not null"`
	// BusinessId           *uuid.UUID     `gorm:"column:business_id;not null"`
	// Business             Business       `gorm:"foreignKey:BusinessId"`
	// BusinessCollectionId *uuid.UUID     `gorm:"column:business_collection_id;not null"`
	// ProvinceId           *uuid.UUID     `gorm:"column:province_id;not null"`
	// MunicipalityId       *uuid.UUID     `gorm:"column:municipality_id;not null"`
	// HighQualityPhoto     string         `gorm:"column:high_quality_photo;not null"`
	// LowQualityPhoto      string         `gorm:"column:low_quality_photo;not null"`
	// Thumbnail            string         `gorm:"column:thumbnail;not null"`
	// BlurHash             string         `gorm:"column:blurhash;not null"`
	// Cursor               int32          `gorm:"column:cursor"`
	// CreateTime           time.Time      `gorm:"column:create_time;not null"`
	// UpdateTime           time.Time      `gorm:"column:update_time;not null"`
	// DeleteTime           gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *itemDatasource) SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64) (*[]entity.ItemBusiness, error) {
	var items []entity.ItemBusiness
	var result *gorm.DB
	where := "%" + name + "%"
	if municipalityNotEqual {
		result = tx.Limit(int(limit+1)).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_id != ? AND province_id = ? AND cursor > ?", where, municipalityId, provinceId, cursor).Order("cursor asc").Find(&items)
	} else {
		result = tx.Limit(int(limit+1)).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND municipality_id = ? AND province_id = ? AND cursor > ?", where, municipalityId, provinceId, cursor).Order("cursor asc").Find(&items)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (i *itemDatasource) SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string) (*[]entity.ItemBusiness, error) {
	var items []entity.ItemBusiness
	var result *gorm.DB
	where := "%" + name + "%"
	result = tx.Limit(10).Where("LOWER(unaccent(item.name)) LIKE LOWER(unaccent(?)) AND cursor > ? AND business_id = ?", where, cursor, businessId).Order("cursor asc").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (v *itemDatasource) UpdateItem(tx *gorm.DB, where *entity.ItemBusiness, data *entity.ItemBusiness) (*entity.ItemBusiness, error) {
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

func (v *itemDatasource) UpdateItems(tx *gorm.DB, data *[]entity.ItemBusiness) (*[]entity.ItemBusiness, error) {
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
