package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UpdateOrderRequest struct {
	Id          uuid.UUID
	Status string
	Metadata    *metadata.MD
}
