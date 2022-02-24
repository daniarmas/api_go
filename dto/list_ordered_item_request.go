package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type ListOrderedItemRequest struct {
	OrderFk  uuid.UUID
	Metadata *metadata.MD
}
