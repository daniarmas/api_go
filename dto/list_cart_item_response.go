package dto

import (
	"time"

	"github.com/daniarmas/api_go/models"
)

type ListCartItemResponse struct {
	CartItems []models.CartItemAndItem
	NextPage  time.Time
}
