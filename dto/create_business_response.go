package dto

import "github.com/daniarmas/api_go/models"

type CreateBusinessResponse struct {
	Business                                     *models.Business
	UnionBusinessAndMunicipalityWithMunicipality *[]models.UnionBusinessAndMunicipalityWithMunicipality
}
