package dto

import "github.com/google/uuid"

type GetBusinessProvinceAndMunicipality struct {
	MunicipalityId *uuid.UUID
	ProvinceId     *uuid.UUID
}
