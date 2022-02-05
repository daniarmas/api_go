package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type ItemQuery interface {
	GetItem(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error)
	ListItem(tx *gorm.DB, where *models.Item, cursor time.Time) (*[]models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error)
	UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error)
}

type itemQuery struct{}

func (i *itemQuery) ListItem(tx *gorm.DB, where *models.Item, cursor time.Time) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().ListItem(tx, where, cursor)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemQuery) GetItem(tx *gorm.DB, id string, point ewkb.Point) (*models.ItemBusiness, error) {
	result, err := Datasource.NewItemDatasource().GetItem(tx, id, point)
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

func (i *itemQuery) SearchItem(tx *gorm.DB, name string, provinceFk string, municipalityFk string, cursor int64, municipalityNotEqual bool, limit int64) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().SearchItem(tx, name, provinceFk, municipalityFk, cursor, municipalityNotEqual, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
