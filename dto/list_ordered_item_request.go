package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type ListOrderedItemRequest struct {
	OrderId  uuid.UUID
	Metadata *metadata.MD
}
