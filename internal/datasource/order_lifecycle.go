package datasource

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type OrderLifecycleDatasource interface {
	CreateOrderLifecycle(tx *gorm.DB, data *entity.OrderLifecycle) (*entity.OrderLifecycle, error)
}

type orderLifecycleDatasource struct{}

func (i *orderLifecycleDatasource) CreateOrderLifecycle(tx *gorm.DB, data *entity.OrderLifecycle) (*entity.OrderLifecycle, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
