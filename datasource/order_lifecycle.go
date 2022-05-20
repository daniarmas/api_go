package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderLifecycleDatasource interface {
	CreateOrderLifecycle(tx *gorm.DB, data *models.OrderLifecycle) (*models.OrderLifecycle, error)
}

type orderLifecycleDatasource struct{}

func (i *orderLifecycleDatasource) CreateOrderLifecycle(tx *gorm.DB, data *models.OrderLifecycle) (*models.OrderLifecycle, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
