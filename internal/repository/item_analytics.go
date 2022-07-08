package repository

import (
	"database/sql"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type ItemAnalyticsRepository interface {
	CreateItemAnalytics(tx *sql.Tx, data *[]entity.ItemAnalytics) (*[]entity.ItemAnalytics, error)
	GetItemAnalytics(tx *gorm.DB, where *entity.ItemAnalytics, fields *[]string) (*entity.ItemAnalytics, error)
	ListItemAnalytics(tx *gorm.DB, where *entity.ItemAnalytics, fields *[]string) (*[]entity.ItemAnalytics, error)
}

type itemAnalyticsRepository struct{}

func (i *itemAnalyticsRepository) CreateItemAnalytics(tx *sql.Tx, data *[]entity.ItemAnalytics) (*[]entity.ItemAnalytics, error) {
	res, err := Datasource.NewItemAnalyticsDatasource().CreateItemAnalytics(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemAnalyticsRepository) ListItemAnalytics(tx *gorm.DB, where *entity.ItemAnalytics, fields *[]string) (*[]entity.ItemAnalytics, error) {
	result, err := Datasource.NewItemAnalyticsDatasource().ListItemAnalytics(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemAnalyticsRepository) GetItemAnalytics(tx *gorm.DB, where *entity.ItemAnalytics, fields *[]string) (*entity.ItemAnalytics, error) {
	result, err := Datasource.NewItemAnalyticsDatasource().GetItemAnalytics(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
