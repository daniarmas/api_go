package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessAnalyticsRepository interface {
	CreateBusinessAnalytics(tx *gorm.DB, data *[]models.BusinessAnalytics) (*[]models.BusinessAnalytics, error)
	GetBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*models.BusinessAnalytics, error)
	ListBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*[]models.BusinessAnalytics, error)
}

type businessAnalyticsRepository struct{}

func (i *businessAnalyticsRepository) CreateBusinessAnalytics(tx *gorm.DB, data *[]models.BusinessAnalytics) (*[]models.BusinessAnalytics, error) {
	res, err := Datasource.NewBusinessAnalyticsDatasource().CreateBusinessAnalytics(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessAnalyticsRepository) ListBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*[]models.BusinessAnalytics, error) {
	result, err := Datasource.NewBusinessAnalyticsDatasource().ListBusinessAnalytics(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessAnalyticsRepository) GetBusinessAnalytics(tx *gorm.DB, where *models.BusinessAnalytics, fields *[]string) (*models.BusinessAnalytics, error) {
	result, err := Datasource.NewBusinessAnalyticsDatasource().GetBusinessAnalytics(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
