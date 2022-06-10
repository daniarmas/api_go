package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessAnalyticsDatasource interface {
	CreateBusinessAnalytics(tx *gorm.DB, data *models.BusinessAnalytics) (*models.BusinessAnalytics, error)
	GetBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*models.BusinessAnalytics, error)
	ListBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*[]models.BusinessAnalytics, error)
}

type businessAnalyticsDatasource struct{}

func (i *businessAnalyticsDatasource) CreateBusinessAnalytics(tx *gorm.DB, data *models.BusinessAnalytics) (*models.BusinessAnalytics, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *businessAnalyticsDatasource) GetBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*models.BusinessAnalytics, error) {
	var res *models.BusinessAnalytics
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

func (i *businessAnalyticsDatasource) ListBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*[]models.BusinessAnalytics, error) {
	var res []models.BusinessAnalytics
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
