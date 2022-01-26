package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ItemCategoryDatasource interface {
	GetItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*models.BusinessItemCategory, error)
	ListItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*[]models.BusinessItemCategory, error)
}

type itemCategoryDatasource struct{}

func (i *itemCategoryDatasource) ListItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*[]models.BusinessItemCategory, error) {
	var itemsCategory []models.BusinessItemCategory
	result := tx.Where(where).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}

func (i *itemCategoryDatasource) GetItemCategory(tx *gorm.DB, where *models.BusinessItemCategory) (*models.BusinessItemCategory, error) {
	var businessItemCategory *models.BusinessItemCategory
	result := tx.Where(where).Take(&businessItemCategory)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessItemCategory, nil
}
