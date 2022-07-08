package repository

import (
	"database/sql"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessAnalyticsRepository interface {
	CreateBusinessAnalytics(tx *sql.Tx, data *[]entity.BusinessAnalytics) (*[]entity.BusinessAnalytics, error)
	GetBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*entity.BusinessAnalytics, error)
	ListBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*[]entity.BusinessAnalytics, error)
}

type businessAnalyticsRepository struct{}

func (i *businessAnalyticsRepository) CreateBusinessAnalytics(tx *sql.Tx, data *[]entity.BusinessAnalytics) (*[]entity.BusinessAnalytics, error) {
	res, err := Datasource.NewBusinessAnalyticsDatasource().CreateBusinessAnalytics(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessAnalyticsRepository) ListBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*[]entity.BusinessAnalytics, error) {
	result, err := Datasource.NewBusinessAnalyticsDatasource().ListBusinessAnalytics(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessAnalyticsRepository) GetBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*entity.BusinessAnalytics, error) {
	result, err := Datasource.NewBusinessAnalyticsDatasource().GetBusinessAnalytics(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
