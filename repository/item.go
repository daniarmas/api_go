package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type ItemQuery interface {
	GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error)
	GetItemWithLocation(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error)
	ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error)
	ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error)
	CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error)
	SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]models.Item, error)
	UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error)
	UpdateItems(tx *gorm.DB, data *[]models.Item) (*[]models.Item, error)
	DeleteItem(tx *gorm.DB, where *models.Item) error
}

type itemQuery struct{}

func (v *itemQuery) DeleteItem(tx *gorm.DB, where *models.Item) error {
	err := Datasource.NewItemDatasource().DeleteItem(tx, where)
	if err != nil {
		return err
	}
	return nil
}

func (v *itemQuery) CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error) {
	res, err := Datasource.NewItemDatasource().CreateItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemQuery) ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().ListItem(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().ListItemInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error) {
	result, err := Datasource.NewItemDatasource().GetItem(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) GetItemWithLocation(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error) {
	result, err := Datasource.NewItemDatasource().GetItemWithLocation(tx, id, point)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error) {
	result, err := Datasource.NewItemDatasource().UpdateItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) UpdateItems(tx *gorm.DB, data *[]models.Item) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().UpdateItems(tx, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().SearchItem(tx, name, provinceId, municipalityId, cursor, municipalityNotEqual, limit, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().SearchItemByBusiness(tx, name, cursor, businessId, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
