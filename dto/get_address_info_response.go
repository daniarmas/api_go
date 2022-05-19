package dto

import "github.com/google/uuid"

type GetAddressInfoResponse struct {
	MunicipalityId           *uuid.UUID
	MunicipalityName         string
	ProvinceNameAbbreviation string
	ProvinceId               *uuid.UUID
	ProvinceName             string
}
