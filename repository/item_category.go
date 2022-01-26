package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ItemCategoryQuery interface {
	GetItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*models.BusinessItemCategory, error)
	ListItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*[]models.BusinessItemCategory, error)
}

type itemCategoryQuery struct{}

func (i *itemCategoryQuery) ListItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*[]models.BusinessItemCategory, error) {
	result, err := Datasource.NewItemCategoryDatasource().ListItemCategory(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemCategoryQuery) GetItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*models.BusinessItemCategory, error) {
	result, err := Datasource.NewItemCategoryDatasource().GetItemCategory(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
