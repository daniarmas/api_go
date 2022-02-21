package dto

import (
	"time"

	"google.golang.org/grpc/metadata"
)

type ListOrderRequest struct {
	NextPage time.Time
	Metadata *metadata.MD
}
