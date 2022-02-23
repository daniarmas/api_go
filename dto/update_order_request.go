package dto

import "github.com/google/uuid"

type UpdateOrderRequest struct {
	Id          uuid.UUID
	OrderStatus string
}
