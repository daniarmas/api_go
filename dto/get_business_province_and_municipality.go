package dto

import "github.com/google/uuid"

type GetBusinessProvinceAndMunicipality struct {
	MunicipalityFk uuid.UUID
	ProvinceFk     uuid.UUID
}
