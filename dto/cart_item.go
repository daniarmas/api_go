package dto

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID                   *uuid.UUID
	Name                 string
	Price                string
	Quantity             int32
	ItemId               *uuid.UUID
	UserId               *uuid.UUID
	AuthorizationTokenId *uuid.UUID
	Thumbnail            string
	ThumbnailBlurHash    string
	Cursor               int32
	CreateTime           time.Time
	UpdateTime           time.Time
}
