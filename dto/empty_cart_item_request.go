package dto

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type EmptyCartItemRequest struct {
	Metadata       *metadata.MD
	MunicipalityFk uuid.UUID
}
