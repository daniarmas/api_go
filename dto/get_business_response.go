package dto

import "github.com/daniarmas/api_go/models"

type GetBusinessResponse struct {
	Business     *models.Business
	ItemCategory *[]models.BusinessItemCategory
}
