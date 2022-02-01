package dto

import "github.com/daniarmas/api_go/models"

type ListCartItemResponse struct {
	CartItems []models.CartItemAndItem
	NextPage  int32
}
