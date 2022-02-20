package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	ListItem(tx *gorm.DB, where *models.Item, cursor time.Time) (*[]models.Item, error)
}

type orderRepository struct{}

func (i *orderRepository) ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error) {
	result, err := Datasource.NewOrderDatasource().ListOrder(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
