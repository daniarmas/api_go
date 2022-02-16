package dto

import (
	"time"

	"github.com/daniarmas/api_go/models"
)

type ListItemResponse struct {
	Items    []models.Item
	NextPage time.Time
}
