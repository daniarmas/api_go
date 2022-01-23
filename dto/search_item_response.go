package dto

import "github.com/daniarmas/api_go/models"

type SearchItemResponse struct {
	Items                  *[]models.Item
	SearchMunicipalityType string
	NextPage               int32
}
