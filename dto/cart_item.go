package dto

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID                   uuid.UUID
	Name                 string
	Price                float64
	Quantity             int32
	ItemFk               uuid.UUID
	UserFk               uuid.UUID
	AuthorizationTokenFk uuid.UUID
	Thumbnail            string
	ThumbnailBlurHash    string
	Cursor               int32
	CreateTime           time.Time
	UpdateTime           time.Time
}
