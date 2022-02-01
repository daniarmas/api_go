package dto

import "google.golang.org/grpc/metadata"

type ListCartItemRequest struct {
	Metadata *metadata.MD
	NextPage int32
}
