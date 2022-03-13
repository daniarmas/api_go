package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

type OrderLifecycleDatasource interface {
	CreateOrder(tx *gorm.DB, data *models.OrderLifecycle) (*models.OrderLifecycle, error)
}

type orderLifecycleDatasource struct{}

func (i *orderLifecycleDatasource) CreateOrder(tx *gorm.DB, data *models.OrderLifecycle) (*models.OrderLifecycle, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
