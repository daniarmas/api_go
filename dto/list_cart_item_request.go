package dto

import (
	"time"

	"google.golang.org/grpc/metadata"
)

type ListCartItemRequest struct {
	Metadata *metadata.MD
	NextPage time.Time
}
