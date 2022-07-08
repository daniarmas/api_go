package repository

import (
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetItem(tx *gorm.DB, where *entity.Item, fields *[]string) (*entity.Item, error)
	ListItem(tx *gorm.DB, where *entity.Item, cursor time.Time, fields *[]string) (*[]entity.Item, error)
	ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.Item, error)
	CreateItem(tx *gorm.DB, data *entity.Item) (*entity.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]entity.Item, error)
	SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]entity.Item, error)
	UpdateItem(tx *gorm.DB, where *entity.Item, data *entity.Item) (*entity.Item, error)
	UpdateItems(tx *gorm.DB, data *[]entity.Item) (*[]entity.Item, error)
	DeleteItem(tx *gorm.DB, where *entity.Item) error
}

type itemRepository struct{}

func (v *itemRepository) DeleteItem(tx *gorm.DB, where *entity.Item) error {
	err := Datasource.NewItemDatasource().DeleteItem(tx, where)
	if err != nil {
		return err
	}
	return nil
}

func (v *itemRepository) CreateItem(tx *gorm.DB, data *entity.Item) (*entity.Item, error) {
	res, err := Datasource.NewItemDatasource().CreateItem(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) ListItem(tx *gorm.DB, where *entity.Item, cursor time.Time, fields *[]string) (*[]entity.Item, error) {
	res, err := Datasource.NewItemDatasource().ListItem(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.Item, error) {
	res, err := Datasource.NewItemDatasource().ListItemInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) GetItem(tx *gorm.DB, where *entity.Item, fields *[]string) (*entity.Item, error) {
	res, err := Datasource.NewItemDatasource().GetItem(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) UpdateItem(tx *gorm.DB, where *entity.Item, data *entity.Item) (*entity.Item, error) {
	res, err := Datasource.NewItemDatasource().UpdateItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) UpdateItems(tx *gorm.DB, data *[]entity.Item) (*[]entity.Item, error) {
	res, err := Datasource.NewItemDatasource().UpdateItems(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]entity.Item, error) {
	res, err := Datasource.NewItemDatasource().SearchItem(tx, name, provinceId, municipalityId, cursor, municipalityNotEqual, limit, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]entity.Item, error) {
	res, err := Datasource.NewItemDatasource().SearchItemByBusiness(tx, name, cursor, businessId, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
