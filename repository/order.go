package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error)
	ListOrderWithBusiness(tx *gorm.DB, where *models.OrderBusiness) (*[]models.OrderBusiness, error)
	CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error)
	UpdateOrder(tx *gorm.DB, where *models.Order, data *models.Order) (*models.Order, error)
}

type orderRepository struct{}

func (i *orderRepository) ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error) {
	result, err := Datasource.NewOrderDatasource().ListOrder(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *orderRepository) ListOrderWithBusiness(tx *gorm.DB, where *models.OrderBusiness) (*[]models.OrderBusiness, error) {
	result, err := Datasource.NewOrderDatasource().ListOrderWithBusiness(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *orderRepository) CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error) {
	res, err := Datasource.NewOrderDatasource().CreateOrder(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderRepository) UpdateOrder(tx *gorm.DB, where *models.Order, data *models.Order) (*models.Order, error) {
	result, err := Datasource.NewOrderDatasource().UpdateOrder(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
