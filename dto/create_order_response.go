package dto

import (
	"time"

	"github.com/daniarmas/api_go/models"
)

type CreateOrderResponse struct {
	OrderedItems *[]models.Order
	NextPage     time.Time
}
