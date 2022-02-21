package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderDatasource interface {
	ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error)
	CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error)
}

type orderDatasource struct{}

func (i *orderDatasource) CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *orderDatasource) ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error) {
	var order []models.Order
	result := tx.Limit(11).Select("id, status, delivery_type, residence_type, price, building_number, house_number, business_fk, user_fk, device_fk, app_version, delivery_date, create_time, update_time, delete_time, ST_AsEWKB(coordinates) AS coordinates").Where("user_fk = ? AND create_time < ?", where.UserFk, where.CreateTime).Order("create_time desc").Find(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
