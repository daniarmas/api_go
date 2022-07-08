package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type OrderLifecycleRepository interface {
	CreateOrderLifecycle(tx *gorm.DB, data *entity.OrderLifecycle) (*entity.OrderLifecycle, error)
}

type orderLifecycleRepository struct{}

func (v *orderLifecycleRepository) CreateOrderLifecycle(tx *gorm.DB, data *entity.OrderLifecycle) (*entity.OrderLifecycle, error) {
	res, err := Datasource.NewOrderLifecycleDatasource().CreateOrderLifecycle(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
