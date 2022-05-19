package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type DeleteItemRequest struct {
	Metadata *metadata.MD
	ItemId   *uuid.UUID
}
