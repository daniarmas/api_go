package dto

import "github.com/daniarmas/api_go/models"

type GetBusinessResponse struct {
	Business            *models.Business
	BusinessCollections *[]models.BusinessCollection
}
