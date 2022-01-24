package dto

import "github.com/daniarmas/api_go/models"

type ListItemResponse struct {
	Items    []models.Item
	NextPage int32
}
