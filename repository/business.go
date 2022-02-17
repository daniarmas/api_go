package repository

import (
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type BusinessQuery interface {
	Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error)
	GetBusiness(tx *gorm.DB, where *models.Business) (*models.Business, error)
	GetBusinessProvinceAndMunicipality(tx *gorm.DB, businessFk uuid.UUID) (*dto.GetBusinessProvinceAndMunicipality, error)
}

type businessQuery struct{}

func (b *businessQuery) Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().Feed(tx, coordinates, limit, provinceFk, municipalityFk, cursor, municipalityNotEqual, homeDelivery, toPickUp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessQuery) GetBusiness(tx *gorm.DB, where *models.Business) (*models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusiness(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessQuery) GetBusinessProvinceAndMunicipality(tx *gorm.DB, businessFk uuid.UUID) (*dto.GetBusinessProvinceAndMunicipality, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusinessProvinceAndMunicipality(tx, businessFk)
	if err != nil {
		return nil, err
	}
	return result, nil
}
