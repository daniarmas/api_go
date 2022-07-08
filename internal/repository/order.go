package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	ListOrder(tx *gorm.DB, where *entity.Order, fields *[]string) (*[]entity.Order, error)
	ListOrderWithBusiness(tx *gorm.DB, where *entity.OrderBusiness) (*[]entity.OrderBusiness, error)
	CreateOrder(tx *gorm.DB, data *entity.Order) (*entity.Order, error)
	UpdateOrder(tx *gorm.DB, where *entity.Order, data *entity.Order) (*entity.Order, error)
	GetOrder(tx *gorm.DB, where *entity.Order) (*entity.Order, error)
}

type orderRepository struct{}

func (i *orderRepository) ListOrder(tx *gorm.DB, where *entity.Order, fields *[]string) (*[]entity.Order, error) {
	result, err := Datasource.NewOrderDatasource().ListOrder(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *orderRepository) GetOrder(tx *gorm.DB, where *entity.Order) (*entity.Order, error) {
	res, err := Datasource.NewOrderDatasource().GetOrder(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderRepository) ListOrderWithBusiness(tx *gorm.DB, where *entity.OrderBusiness) (*[]entity.OrderBusiness, error) {
	result, err := Datasource.NewOrderDatasource().ListOrderWithBusiness(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *orderRepository) CreateOrder(tx *gorm.DB, data *entity.Order) (*entity.Order, error) {
	res, err := Datasource.NewOrderDatasource().CreateOrder(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *orderRepository) UpdateOrder(tx *gorm.DB, where *entity.Order, data *entity.Order) (*entity.Order, error) {
	result, err := Datasource.NewOrderDatasource().UpdateOrder(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
