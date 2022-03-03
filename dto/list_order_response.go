package dto

import (
	"time"

	"github.com/daniarmas/api_go/models"
)

type ListOrderResponse struct {
	NextPage time.Time
	Orders   *[]models.OrderBusiness
}
