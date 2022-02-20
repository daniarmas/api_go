package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type OrderDatasource interface {
	ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error)
}

type orderDatasource struct{}

func (i *orderDatasource) ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error) {
	var order []models.Order
	result := tx.Limit(11).Where("user_fk = ? AND create_time < ?", where.UserFk, where.CreateTime).Order("create_time desc").Find(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
