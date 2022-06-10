package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ItemAnalyticsDatasource interface {
	CreateItemAnalytics(tx *gorm.DB, data *[]models.ItemAnalytics) (*[]models.ItemAnalytics, error)
	GetItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*models.ItemAnalytics, error)
	ListItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*[]models.ItemAnalytics, error)
}

type itemAnalyticsDatasource struct{}

func (i *itemAnalyticsDatasource) CreateItemAnalytics(tx *gorm.DB, data *[]models.ItemAnalytics) (*[]models.ItemAnalytics, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *itemAnalyticsDatasource) GetItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*models.ItemAnalytics, error) {
	var res *models.ItemAnalytics
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

func (i *itemAnalyticsDatasource) ListItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*[]models.ItemAnalytics, error) {
	var res []models.ItemAnalytics
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
