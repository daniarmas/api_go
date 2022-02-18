package dto

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type ListCartItemRequest struct {
	Metadata       *metadata.MD
	NextPage       time.Time
	MunicipalityFk uuid.UUID
}
